package rapidroot

import (
	"strings"
)

type node struct {
	pathSegment           string
	dynamicValue          string
	handler               HandlerFunc
	groupMiddleware       []Middleware
	currentNodeMiddleware []Middleware
	children              []*node
	isDynamic             bool
}

func newNode() *node {
	return &node{
		children:              make([]*node, 0),
		currentNodeMiddleware: make([]Middleware, 0),
		groupMiddleware:       make([]Middleware, 0),
	}
}

func (n *node) freeMiddlewaresMemory() {
	n.currentNodeMiddleware = nil
	n.groupMiddleware = nil
}

func (n *node) child(pathSegment string) *node {
	for _, child := range n.children {
		if child.pathSegment == pathSegment {
			return child
		}
	}
	return nil
}

func (n *node) addChild(child *node) {
	n.children = append(n.children, child)
}

func addRoute(path string, root *node, handler HandlerFunc) *node {
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

		if isLastSegment(i, segments) {
			currentNode.handler = handler
		}
	}

	return currentNode
}

func getNode(path string, root *node, req *Request) *node {
	segments := strings.Split(path, "/")
	currentNode := root

	for i, segment := range segments {
		childNode := currentNode.child(segment)

		if childNode == nil {
			dynamicChild := findDynamicChild(currentNode.children, segments[i+1:])
			if dynamicChild == nil {
				// If no matching child is found and no dynamic segment is available,
				// it means the route is not defined. Handle this case accordingly.
				return nil
			}

			currentNode = dynamicChild
			if req != nil {
				req.SetValue(dynamicChild.pathSegment, segment)
			}
		} else {
			currentNode = childNode

			if currentNode.isDynamic && req != nil {
				req.SetValue(currentNode.pathSegment, segment)
			}
		}

		if i == len(segments)-1 {
			return currentNode
		}
	}

	return currentNode
}

func findDynamicChild(children []*node, remainingSegments []string) *node {
	for _, dynamicChild := range children {
		if dynamicChild.isDynamic {
			nextNode := getNode(strings.Join(remainingSegments, "/"), dynamicChild, nil)
			if nextNode != nil {
				return dynamicChild
			}
		}
	}

	return nil
}
