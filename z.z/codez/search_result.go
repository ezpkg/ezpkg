package codez

import (
	"go/ast"
)

type SearchResult struct {
	Items []SearchResultItem
}

type SearchResultItem struct {
	Node ast.Node

	PatternID int
}
