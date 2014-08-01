# gogl

[![Build Status](https://travis-ci.org/sdboyer/gogl.png?branch=master)](https://travis-ci.org/sdboyer/gogl)
[![Coverage Status](https://coveralls.io/repos/sdboyer/gogl/badge.png?branch=master)](https://coveralls.io/r/sdboyer/gogl?branch=master)

gogl is a graph library in Go. Its goal is to provide simple, unifying interfaces and implementations of graph algorithms and datastructures that can scale from small graphs to very large graphs. The latter case is, as yet, untested!

gogl is based on the premise that working with graphs can be [decomplected](http://www.infoq.com/presentations/Simple-Made-Easy) by focusing primarily on the natural constraints established in graph theory.

There's still a lot to do - gogl is still firming up significant aspects of how its API works.

## Principles

Graph systems are often big, complicated affairs. gogl tries to be not that. These are the operant principles:

1. Simplicity: fully and correctly modeling graph theoretic concepts in idiomatic Go.
1. Performance: be as fast as design constraints and known-best algorithms allow.
1. Extensibility: expect others to run gogl's graph datastructures through their own algorithms, , and gogl's algorithms with their graph implementations.
1. Functional: orient towards transforms, functors, and streams; achieve other styles through layering.
1. Flexibility: Be unopinionated about vertices, and minimally opinionated about edges.
1. Correctness: Utilize [commonly accepted graph terminology](http://en.wikipedia.org/wiki/Glossary_of_graph_theory) where possible, and adhere to its meaning.

The first and last points are key - names in gogl are carefully chosen, with the hope that they can guide intuition when stricter rules (e.g., the type system) become ambiguous. The [godoc](https://godoc.org/github.com/sdboyer/gogl) generally takes care to detail these subtleties. But godoc is a reference, not a tutorial.

## Quickstart

Getting started with gogl is simple: create a graph object, add your data, and off you go.

```go
package main

import (
	"fmt"
	"github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/dfs"
)

func main() {
	// gogl uses a builder to specify the kind of graph you want.
	graph := gogl.G().
		// The graph should be mutable. Default is immutable.
		Mutable().
		// The graph should have directed edges (arcs). Default is undirected.
		Directed().
		// The graph's edges are plain - no labels, weights, etc. This is the default.
		Basic().
		// No loops or parallel edges. This is the default.
		SimpleGraph().
		// gogl.AdjacencyList picks and returns an adjacency list-based graph, based on the spec.
		Create(gogl.AdjacencyList).
		// The builder always returns a Graph; type assert to get access to add/remove methods.
		(gogl.MutableGraph)

	// Adds two basic edges. Of course, this adds the vertices, too.
	graph.AddEdges(gogl.NewEdge("foo", "bar"), gogl.NewEdge("bar", "baz"))

	// gogl's core iteration concept is built on injected functions (VertexLambda or
	// EdgeLambda). Here, a VertexLambda is called once per vertex in the graph;
	// the return value determines whether traversal continues.
	graph.EachVertex(func(v gogl.Vertex) (terminate bool) {
		fmt.Println(v) // Probably "foo\nbar\nbaz", but ordering is not guaranteed.
		return // returns false, so iteration continues
	})

	// gogl refers to these sorts of iterating methods as enumerators. There are four
	// such methods on undirected graphs, and two more on directed graphs.

	// If you know you need the full result set, gogl provides functors to collect enumerations
	// into slices. This makes ranging easy.
	var vertices []gogl.Vertex = gogl.CollectVertices(graph)
	for _, v := range vertices {
		fmt.Println(v) // same as with EachVertex().
	}

	// The pattern is the same with edge enumeration. These two have the same output:
	graph.EachEdge(func(e gogl.Edge) (terminate bool) {
		fmt.Println(e) // Probably "{foo bar}\n{bar baz}". Again, ordering is not guaranteed.
		return
	})
	for _, e := range gogl.CollectEdges(graph) {
		fmt.Println(e)
	}

	// gogl's algorithms all rely on these enumerators to do their work. Here, we use
	// a depth-first topological sort algorithm to produce a slice of vertices.
	var tsl []gogl.Vertex
	tsl, err := dfs.Toposort(graph, "foo")
	if err == nil {
		fmt.Println(tsl) // [baz bar foo]
	}
}
```

## Enumerators

Enumerators are the primary means by which gogl graphs are expressed. As shown in the Quickstart section, they are methods on graph datastructures that receive a 'lambda', and call that lambda once per element (Vertex or Edge) that is found as the method traverses the graph. There are four enumerators for gogl's undirected graphs, and two additional ones for directed graphs.

Given the following graph:
![Base graph](doc/base.dot.png)

Which could be created as follows:
```go
func main() {
	graph := gogl.G().Mutable().Directed().Create(gogl.AdjacencyList)

	graph.AddEdges([]gogl.Edge{
		NewEdge("a", "b"),
		NewEdge("b", "c"),
		NewEdge("a", "c"),
		NewEdge("a", "c"),
		NewEdge("d", "a"),
		NewEdge("d", "e"),
	})

	// 'f' is a vertex isolate.
	graph.EnsureVertex("f")
}
```

Calling `EachVertex()` on this graph will call the injected lambda six times, once for each of the contained vertices (marked in blue). It's important to remember that gogl makes no guarantees as to the order of the elements.
![EachVertex()](doc/ev.dot.png)

Calling `EachEdge()` will call the injected lambda six times, once for each of the contained edges:
![EachEdge()](doc/ee.dot.png)

## Gotchas
