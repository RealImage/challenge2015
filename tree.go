package main

var tree map[string][]string

func AddChildren(parName string, children []string) {
	tree[parName] = children
}
