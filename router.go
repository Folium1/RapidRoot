package rapidRoot

import (
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	root       *node
	routesList strings.Builder
}

// NewRouter returns a new router instance with default configuration.
func NewRouter() *Router {
	r := &Router{
		root: newNode(),
	}
	r.routesList.WriteString("\n------------Handlers--------------\n\n")
	return r
}

type node struct {
	pathSegment     string
	method          string
	handler         HandlerFunc
	staticSegments  map[string]*node
	dynamicSegments map[string]*node
	middlewares     []Middleware
}

func newNode() *node {
	return &node{
		pathSegment:     "",
		method:          "",
		handler:         nil,
		staticSegments:  make(map[string]*node),
		dynamicSegments: make(map[string]*node),
		middlewares:     make([]Middleware, 0),
	}
}

func (r *Router) handle(method, path string, handler HandlerFunc) {
	r.routesList.WriteString(fmt.Sprintf("%s %s %s\n", method, path, getFunctionName(handler)))
	if handler == nil {
		log.fatal(fmt.Errorf("Nil handlers are not allowed | %s %s %s\n", method, path, "nil"))
	}
	currentNode := r.root
	pathSegments := strings.Split(path, "/")

	for _, pathSegment := range pathSegments {
		if currentNode.staticSegments == nil {
			currentNode.staticSegments = make(map[string]*node)
		}

		if !strings.HasPrefix(pathSegment, "$") {
			// handle static segments
			if child, exists := currentNode.staticSegments[pathSegment]; exists {
				currentNode = child
			} else {
				newNode := newNode()
				currentNode.staticSegments[pathSegment] = newNode
				currentNode = newNode
			}
		} else {
			// handle wildcard
			dynamicSegment := pathSegment[1:]
			if _, exists := currentNode.staticSegments["$"]; !exists {
				newNode := newNode()
				currentNode.staticSegments["$"] = newNode
				currentNode = newNode
				currentNode.pathSegment = dynamicSegment
			} else {
				currentNode = currentNode.staticSegments["$"]
			}
		}
	}

	currentNode.method = method
	currentNode.handler = handler
}

func (r *Router) getHandler(method, path string, req *Request) HandlerFunc {
	currentNode := r.root
	pathSegments := strings.Split(path, "/")

	for _, pathSegment := range pathSegments {
		if currentNode == nil {
			return nil
		}
		if child, isWild := currentNode.staticSegments["$"]; isWild {
			currentNode = child
			req.SetValue(currentNode.pathSegment, pathSegment)
		} else if staticChild, ok := currentNode.staticSegments[pathSegment]; ok {
			currentNode = staticChild
		} else {
			return nil
		}
	}

	if currentNode.method == method {
		return currentNode.handler
	}

	return nil
}

type Middleware func(HandlerFunc) HandlerFunc

// MIDDLEWARE adds middleware to the route and all its children.
// If path is "/" then the middleware will be applied to all routes
// middlewares will execute in the order they are in the route
// for example:
//
//	r.MIDDLEWARE("/", middleware1, middleware2)
//	r.GET("/", handler)
//	r.GET("/test", handler)
//	r.Run(":8080")
//
// middleware1 and middleware2 will be applied to all routes. And middleware1 will be executed before middleware2.
// If you want to apply middleware to a specific route, you can do it like this:
//
//	r.MIDDLEWARE("/test", middleware1, middleware2)
//	r.GET("/", handler)
//	r.GET("/test", handler)
//	r.Run(":8080")
func (r *Router) MIDDLEWARE(path string, middlewares ...Middleware) {
	if path == "/" {
		newRoot := newNode()
		newRoot.staticSegments = r.root.staticSegments
		newRoot.middlewares = middlewares
		return
	}

	currentNode := r.root
	pathSegments := strings.Split(path, "/")

	for _, pathSegment := range pathSegments {
		if currentNode.staticSegments == nil {
			currentNode.staticSegments = make(map[string]*node)
		}

		if !strings.HasPrefix(pathSegment, "$") {
			// handle static segments
			if child, exists := currentNode.staticSegments[pathSegment]; exists {
				currentNode = child
			} else {
				newNode := newNode()
				currentNode.staticSegments[pathSegment] = newNode
				currentNode = newNode
			}
		} else {
			// handle wildcard
			dynamicSegment := pathSegment[1:]
			if _, exists := currentNode.staticSegments["$"]; !exists {
				newNode := newNode()
				currentNode.staticSegments["$"] = newNode
				currentNode = newNode
				currentNode.pathSegment = dynamicSegment
			} else {
				currentNode = currentNode.staticSegments["$"]
			}
		}
	}
	currentNode.middlewares = append(currentNode.middlewares, middlewares...)
}

// Helper method to apply middleware for the root node and all children.
func (r *Router) applyMiddlewareForRoutes() {
	r.traverseAndApplyMiddleware(r.root, r.root.middlewares...)
}

// Recursive method to traverse the route tree and apply middleware.
func (r *Router) traverseAndApplyMiddleware(currentNode *node, middleWares ...Middleware) {
	if currentNode == nil {
		return
	}

	currentNode.middlewares = append(currentNode.middlewares, middleWares...)

	// Apply middlewares to the current node
	for _, middleware := range currentNode.middlewares {
		currentNode.handler = middleware(currentNode.handler)
	}

	// Recursively traverse children with a copy of the middlewares
	for _, child := range currentNode.staticSegments {
		r.traverseAndApplyMiddleware(child, currentNode.middlewares...)
	}
	for _, child := range currentNode.dynamicSegments {
		r.traverseAndApplyMiddleware(child, currentNode.middlewares...)
	}
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.handle(http.MethodGet, cleanPath(path), handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.handle(http.MethodPost, cleanPath(path), handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.handle(http.MethodDelete, cleanPath(path), handler)
}

func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.handle(http.MethodPatch, cleanPath(path), handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) {
	r.handle(http.MethodPut, cleanPath(path), handler)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc) {
	r.handle(http.MethodOptions, cleanPath(path), handler)
}

func (r *Router) HEAD(path string, handler HandlerFunc) {
	r.handle(http.MethodHead, cleanPath(path), handler)
}

func (r *Router) CONNECT(path string, handler HandlerFunc) {
	r.handle(http.MethodConnect, cleanPath(path), handler)
}

func (r *Router) TRACE(path string, handler HandlerFunc) {
	r.handle(http.MethodTrace, cleanPath(path), handler)
}
