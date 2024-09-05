package codez

import (
	"go/ast"
)

type _NodeReplaceFunc func(parent ast.Node, idx int, new ast.Node) error

type NodeX struct {
	ast.Node
	parent   ast.Node
	replacer _NodeReplaceFunc
}

func newNodeX(node ast.Node, parent ast.Node, replacer _NodeReplaceFunc) *NodeX {
	return &NodeX{
		Node:     node,
		parent:   parent,
		replacer: replacer,
	}
}
