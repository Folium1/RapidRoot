package rapidroot

import "strings"

type Middleware func(HandlerFunc) HandlerFunc

func (r *Router) GroupMiddleware(method, path string, middleware ...Middleware) {
	if len(middleware) == 0 {
		return
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
			root.groupMiddleware = append(root.groupMiddleware, middleware...)
		}
	}

	addMiddlewareOrGroupToTree(path, root, middleware, true)
}

func (r *Router) Middleware(method, path string, middleware ...Middleware) {
	root := r.getOrCreateRoot(method)
	node := getNode(path, root, nil)
	if node == nil {
		addMiddlewareOrGroupToTree(cleanPath(path), root, middleware, false)
	} else {
		node.currentNodeMiddleware = append(node.currentNodeMiddleware, middleware...)
	}
}

func addMiddlewareOrGroupToTree(path string, root *node, middleware []Middleware, isGroup bool) {
	segments := strings.Split(path, "/")
	currentNode := root

	for i, segment := range segments {
		childNode := currentNode.child(segment)

		if childNode == nil {
			childNode = newNode()
			childNode.pathSegment = segment
			currentNode.addChild(childNode)
		}

		currentNode = childNode

		if isDynamicSegment(segment) {
			currentNode.isDynamic = true
			currentNode.pathSegment = segment[1:]
		}

		if i == len(segments)-1 {
			if isGroup {
				currentNode.groupMiddleware = append(currentNode.groupMiddleware, middleware...)
			} else {
				currentNode.currentNodeMiddleware = append(currentNode.currentNodeMiddleware, middleware...)
			}
		}
	}
}

func (r *Router) applyMiddlewareForRoutes() {
	// Apply middlewares for the root node
	for _, node := range r.tree {
		if node.handler != nil {
			r.applyMiddlewares(node)
		}
		r.traverseAndApplyMiddleware(node, node.groupMiddleware...)
	}

	// Free memory for the root nodes
	for _, node := range r.tree {
		node.freeMiddlewaresMemory()
	}
}

func (r *Router) applyMiddlewares(currentNode *node) {
	for i := len(currentNode.currentNodeMiddleware) - 1; i >= 0; i-- {
		currentNode.handler = currentNode.currentNodeMiddleware[i](currentNode.handler)
	}

	for _, parentMiddleware := range currentNode.groupMiddleware {
		currentNode.handler = parentMiddleware(currentNode.handler)
	}
}

func (r *Router) traverseAndApplyMiddleware(node *node, parentMiddlewares ...Middleware) {
	if node == nil {
		return
	}

	allMiddlewares := append(parentMiddlewares, node.groupMiddleware...)
	allMiddlewares = append(allMiddlewares, node.currentNodeMiddleware...)

	for _, child := range node.children {
		r.traverseAndApplyMiddleware(child, allMiddlewares...)
	}

	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		node.handler = allMiddlewares[i](node.handler)
	}

	node.freeMiddlewaresMemory()
}
