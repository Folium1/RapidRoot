package rapidroot

import (
	"fmt"
	"net/http"
)

type Router struct {
	tree       map[string]*node
	routesList []byte
}

// NewRouter returns a new router instance with default configuration.
func NewRouter() *Router {
	return &Router{
		tree:       make(map[string]*node),
		routesList: []byte("\n------------Handlers--------------\n\n"),
	}
}

// Helper method to get or create the root node for a specific HTTP method.
func (r *Router) getOrCreateRoot(method string) *node {
	if root, ok := r.tree[method]; ok {
		return root
	}
	root := newNode()
	r.tree[method] = root
	return root
}

func (r *Router) handle(method, path string, handler HandlerFunc) {
	r.routesList = append(r.routesList, []byte(fmt.Sprintf("%s %s %s\n", method, path, getFunctionName(handler)))...)
	if handler == nil {
		log.fatal(fmt.Errorf("nil handlers are not allowed | %s %s %s\n", method, path, "nil"))
	}

	path = cleanPath(path)
	root := r.getOrCreateRoot(method)

	if path == "/" {
		if root.pathSegment != "/" {
			newRoot := newNode()
			newRoot.pathSegment = "/"
			newRoot.children = append(newRoot.children, root)
			r.tree[method] = newRoot
		} else {
			root.handler = handler
		}
		return
	}

	node := getNode(path, root, nil)
	if node == nil {
		node = addRoute(path, root, handler)
	}
	node.handler = handler
}

func (r *Router) getHandler(method, path string, req *Request) HandlerFunc {
	if root, ok := r.tree[method]; ok {
		currentNode := getNode(path, root, req)
		if currentNode == nil {
			return nil
		}
		return currentNode.handler
	}
	return nil
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.handle(http.MethodPost, path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.handle(http.MethodDelete, path, handler)
}

func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.handle(http.MethodPatch, path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) {
	r.handle(http.MethodPut, path, handler)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc) {
	r.handle(http.MethodOptions, path, handler)
}

func (r *Router) HEAD(path string, handler HandlerFunc) {
	r.handle(http.MethodHead, path, handler)
}

func (r *Router) CONNECT(path string, handler HandlerFunc) {
	r.handle(http.MethodConnect, path, handler)
}

func (r *Router) TRACE(path string, handler HandlerFunc) {
	r.handle(http.MethodTrace, path, handler)
}
