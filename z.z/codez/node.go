package codez

import (
	"go/ast"
)

type _NodeReplacer func(parent, new ast.Node)

type NodeX struct {
	ast.Node
	parent   ast.Node
	replacer _NodeReplacer
}

func newNodeX(node ast.Node, parent ast.Node, replacer _NodeReplacer) *NodeX {
	return &NodeX{
		Node:     node,
		parent:   parent,
		replacer: replacer,
	}
}
