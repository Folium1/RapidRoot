package rapidRoot

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

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	responseWrapper := &responseCodeWrapper{w, 0}
	reqStruct := newReq(responseWrapper, req)
	handler := r.getHandler(req.Method, cleanPath(req.URL.Path), reqStruct)
	if handler == nil {
		notFoundHandler(w, req)
		return
	}
	reqStruct.handlerName = getFunctionName(handler)
	done := make(chan struct{})
	go func() {
		handler(reqStruct)
		done <- struct{}{}
		close(done)
	}()
	<-done

	log.logRequest(req.URL.Path, req.Method, responseWrapper.statusCode)
}

// Run starts the HTTP server.
func (r *Router) Run(addr string) {
	if log.output == os.Stdout {
		r.routesList.WriteString("\n----------------------------------\n\n")
		logHandlers(r.routesList.String())
		r.routesList.Reset()
	}
	r.applyMiddlewareForRoutes()
	if err := http.ListenAndServe(addr, r); err != nil {
		log.fatal(fmt.Errorf("%s : %w", "Couldn't start the server, err:", err))
	}
}

// RunWithTLS starts the HTTPS server.
func (r *Router) RunWithTLS(addr string, certFile string, keyFile string) {
	if log.output == os.Stdout {
		r.routesList.WriteString("\n----------------------------------\n\n")
		logHandlers(r.routesList.String())
		r.routesList.Reset()
	}
	r.applyMiddlewareForRoutes()
	if err := http.ListenAndServeTLS(addr, certFile, keyFile, r); err != nil {
		log.fatal(fmt.Errorf("%s : %w", "Couldn't start the server, err:", err))
	}
}
