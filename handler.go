package rapidroot

import (
	"net/http"
)

// HandlerFunc is a function that can be registered to a router to handle HTTP
// requests.
type HandlerFunc func(*Request)

// notFoundHandler return 404 with not found message.
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.logRequest(r.URL.Path, r.Method, http.StatusNotFound)
	http.NotFound(w, r)
}

