package main

import (
	"fmt"
	"strings"
)

type AliasTree map[int]AliasNode

type AliasNode struct {
	Alias string
	Nodes AliasTree
}

func (node AliasNode) GetAlias(indexes ...int) (string, bool) {
	if len(indexes) == 0 {
		return node.Alias, node.Alias != ""
	}

	if node.Nodes == nil {
		return fmt.Sprintf("%s.%s", node.Alias, intsJoin(indexes, ".")), false
	}

	subAlias, ok := node.Nodes.GetAlias(indexes...)
	return fmt.Sprintf("%s.%s", node.Alias, subAlias), ok
}

func (tree AliasTree) GetAlias(indexes ...int) (string, bool) {
	if tree == nil || len(indexes) == 0 {
		return "", false
	}
	head, tail := indexes[0], indexes[1:]
	node, ok := tree[head]
	if !ok {
		return intsJoin(tail, "."), false
	}

	return node.GetAlias(tail...)
}

func intsJoin(ints []int, sep string) string {
	var sb strings.Builder
	for idx, item := range ints {
		if idx != 0 {
			sb.WriteString(sep)
		}
		fmt.Fprint(&sb, item)
	}
	return sb.String()
}
