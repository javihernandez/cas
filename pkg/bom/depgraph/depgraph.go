/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package depgraph

import (
	"sync"

	"github.com/codenotary/cas/pkg/bom/artifact"
)

// need a wrapper around dependency to allow adding links to asset before asset is processed
type GraphNode struct {
	asset   *artifact.Dependency
	depType artifact.DepType
}

type graphLink struct {
	from *GraphNode
	to   *GraphNode
}

type key struct {
	name    string
	version string
}

type Graph struct {
	mutex *sync.Mutex
	Root  *GraphNode
	links []graphLink
	nodes map[key]*GraphNode
}

// NewGraph creates new dependency graph with specified root
func NewGraph(name, version string) Graph {
	root := &GraphNode{}
	nodes := make(map[key]*GraphNode)
	nodes[key{name, version}] = root
	g := Graph{Root: root, nodes: nodes, mutex: &sync.Mutex{}}

	return g
}

// Node returns existing node or creates a new empty one
func (g *Graph) Node(name, version string) *GraphNode {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	node, ok := g.nodes[key{name, version}]
	if ok {
		return node
	}
	node = &GraphNode{depType: artifact.DepTransient}
	g.nodes[key{name, version}] = node
	return node
}

// NewNode creates a new node for the dependency
func (g *Graph) NewNode(name, version string, asset *artifact.Dependency) *GraphNode {
	node := g.Node(name, version)
	if node.asset == nil {
		node.asset = asset
	}
	return node
}

// AddChild adds relation between components. If node is nil, dependency added as direct one
func (g *Graph) AddChild(parent *GraphNode, name, version string, child *artifact.Dependency) *GraphNode {
	node := g.Node(name, version)
	if node.asset == nil {
		node.asset = child // set only if wasn't set before
	}
	// transient by default
	if parent == g.Root {
		node.depType = artifact.DepDirect
	}
	g.links = append(g.links, graphLink{parent, node})

	return node
}

// FlatDeps return flat slice with all dependencies (excluding the root)
func (g *Graph) FlatDeps() []artifact.Dependency {
	res := make([]artifact.Dependency, 0, len(g.nodes))
	for _, node := range g.nodes {
		if node == g.Root || node.asset == nil {
			continue
		}
		node.asset.Type = node.depType
		res = append(res, *node.asset)
	}
	return res
}
