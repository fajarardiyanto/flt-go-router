package tree

import (
	"github.com/fajarardiyanto/flt-go-router/interfaces"
	"strings"
)

type (
	Tree struct {
		Root   *Node
		Routes map[string]*Node
	}

	Node struct {
		Key        string
		Path       string
		Handle     interfaces.Handler
		Depth      int
		Children   map[string]*Node
		IsPattern  bool
		Middleware []interfaces.MiddlewareFunc
	}
)

// NewNode returns a newly initialized Node object tha implements the Node
func NewNode(key string, depth int) *Node {
	return &Node{
		Key:      key,
		Depth:    depth,
		Children: make(map[string]*Node),
	}
}

// NewTree returns a newly initialized Tree object that implements the Tree
func NewTree() *Tree {
	return &Tree{
		Root:   NewNode("/", 1),
		Routes: make(map[string]*Node),
	}
}

// Add use `pattern` 、 Handle 、 Middleware stack as node register to tree
func (t *Tree) Add(pattern string, handle interfaces.Handler, middleware ...interfaces.MiddlewareFunc) {
	var currentNode = t.Root

	if pattern != currentNode.Key {
		pattern = TrimPathPrefix(pattern)
		res := SplitPattern(pattern)
		for _, key := range res {
			node, ok := currentNode.Children[key]
			if !ok {
				node = NewNode(key, currentNode.Depth+1)
				node = NewNode(key, currentNode.Depth+1)
				if len(middleware) > 0 {
					node.Middleware = append(node.Middleware, middleware...)
				}
				currentNode.Children[key] = node
			}
			currentNode = node
		}
	}

	if len(middleware) > 0 && currentNode.Depth == 1 {
		currentNode.Middleware = append(currentNode.Middleware, middleware...)
	}

	currentNode.Handle = handle
	currentNode.IsPattern = true
	currentNode.Path = pattern
}

// Find returns nodes that the request match the route pattern
func (t *Tree) Find(pattern string, isRegex bool) (nodes []*Node) {
	var (
		node  = t.Root
		queue []*Node
	)

	if pattern == node.Path {
		nodes = append(nodes, node)
		return
	}

	if !isRegex {
		pattern = TrimPathPrefix(pattern)
	}

	res := SplitPattern(pattern)

	for _, key := range res {
		child, ok := node.Children[key]

		if !ok && isRegex {
			break
		}

		if !ok && !isRegex {
			return
		}

		if pattern == child.Path && !isRegex {
			nodes = append(nodes, child)
			return
		}
		node = child
	}

	queue = append(queue, node)

	for len(queue) > 0 {
		var queueTemp []*Node
		for _, n := range queue {
			if n.IsPattern {
				nodes = append(nodes, n)
			}

			for _, childNode := range n.Children {
				queueTemp = append(queueTemp, childNode)
			}
		}

		queue = queueTemp
	}
	return
}

// TrimPathPrefix is short for strings.TrimPrefix with param prefix `/`
func TrimPathPrefix(pattern string) string {
	return strings.TrimPrefix(pattern, "/")
}

// SplitPattern is short for strings.Split with param seq `/`
func SplitPattern(pattern string) []string {
	return strings.Split(pattern, "/")
}
