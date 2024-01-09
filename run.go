package rapidroot

import (
	"fmt"
	"net/http"
	"os"
)

type responseCodeWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseCodeWrapper) WriteHeader(statusCode int) {
	if w.statusCode != 0 {
		log.warn("Couldn't change status of the response, it had already been changed")
		return
	}
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseCodeWrapper) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := &responseCodeWrapper{w, 0}
	reqStruct := getRequest(resp, req)
	defer releaseRequest(reqStruct)

	handler := r.getHandler(req.Method, cleanPath(req.URL.Path), reqStruct)
	if handler == nil {
		notFoundHandler(w, req)
		return
	}
	reqStruct.handlerName = getFunctionName(handler)

	done := make(chan struct{})
	go func() {
		handlerWrapper(handler, reqStruct)
		done <- struct{}{}
		close(done)
	}()
	<-done

	log.logRequest(req.URL.Path, req.Method, resp.statusCode)
}

func handlerWrapper(handler HandlerFunc, req *Request) {
	handler(req)
	if req.Writer.(*responseCodeWrapper).statusCode == 0 {
		req.SetStatus(http.StatusOK)
	}
}

// Run starts the HTTP server.
func (r *Router) Run(addr string) {
	r.addRoutesListSeparator()
	r.applyMiddlewareForRoutes()
	if err := http.ListenAndServe(addr, r); err != nil {
		log.fatal(fmt.Errorf("Couldn't start the server: %w", err))
	}
}

// addRoutesListSeparator adds a separator to the routes list for debugging purposes.
func (r *Router) addRoutesListSeparator() {
	if log.output != os.Stdout {
		return
	}
	r.routesList = append(r.routesList, []byte("\n----------------------------------\n\n")...)
	logHandlers(string(r.routesList))
	r.routesList = nil
}

// RunWithTLS starts the HTTPS server.
func (r *Router) RunWithTLS(addr, certFile, keyFile string) {
	r.applyMiddlewareForRoutes()
	err := http.ListenAndServeTLS(addr, certFile, keyFile, r)
	if err != nil {
		log.fatal(fmt.Errorf("Couldn't start the server, err: %w", err))
	}
}
