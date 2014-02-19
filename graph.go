package gogl

import (
	"fmt"
	"reflect"
)

/* Vertex structures */

// As a rule, gogl tries to place as low a requirement on its vertices as
// possible. This is because, from a purely graph theoretic perspective,
// vertices are inert. Boring, even. Graphs are more about the topology, the
// characteristics of the edges connecting the points than the points
// themselves. Your use case cares about the content of your vertices, but gogl
// does not.  Consequently, anything can act as a vertex.
type Vertex interface{}
type VertexList []Vertex

// VertexSet uses maps to express a value-less (empty struct), indexed
// unordered list. See
// https://groups.google.com/forum/#!searchin/golang-nuts/map/golang-nuts/H2cXpwisEUE/1X2FV-rODfIJ
type VertexSet map[Vertex]struct{}

// However, as a practical, internal matter, some approaches to representing
// graphs may benefit significantly (in terms of memory use) from having a uint
// that uniquely identifies each vertex. This is essentially an implementation
// detail, but if your use case merits it (large graph processing is busting
// memory), you can trade off a bit of speed for potentially nontrivial memory
// savings.
// TODO implement a datastructure that does this, and test these outlandish claims!
type IdentifiableVertex interface {
	id() uint64
}

/* Edge structures */

// A graph's behaviors is primarily a product of the constraints and
// capabilities it places on its edges. These constraints and capabilities
// determine whether certain types of operations are possible on the graph, as
// well as the efficiencies for various operations.

// gogl aims to provide a diverse range of graph implementations that can meet
// the varying constraints and implementation needs, but still achieve optimal
// performance given those constraints.

// TODO totally unclear whether or not defining capabilities in a bitfield like
// this will actually help us achieve the goal
const (
	E_DIRECTED, EM_DIRECTED = 1 << iota, 1<<iota - 1
	E_UNDIRECTED, EM_UNDIRECTED
	E_WEIGHTED, EM_WEIGHTED
	E_TYPED, EM_TYPED
	E_SIGNED, EM_SIGNED
	E_LOOPS, EM_LOOPS
	E_MULTIGRAPH, EM_MULTIGRAPH
)

type Edge interface {
	Source() Vertex
	Target() Vertex
	Both() (Vertex, Vertex)
	//  Properties() uint
}

type EdgeList []Edge
type Path []Edge

// BaseEdge is a struct used internally to represent edges and meet the Edge
// interface requirements. It uses the standard notation, (u,v), for vertex
// pairs in an edge.
type BaseEdge struct {
	U Vertex
	V Vertex
}

func (e BaseEdge) Source() Vertex {
	return e.U
}

func (e BaseEdge) Target() Vertex {
	return e.V
}

func (e BaseEdge) Both() (Vertex, Vertex) {
	return e.U, e.V
}

/* Graph structures */

type Graph interface {
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(edge Edge))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
	Order() int
	Size() int
	InDegree(vertex Vertex) (int, bool)
	OutDegree(vertex Vertex) (int, bool)
}
type MutableGraph interface {
	Graph
	EnsureVertex(vertices ...Vertex)
	RemoveVertex(vertices ...Vertex)
	AddEdge(edge Edge) bool
	RemoveEdge(edge Edge)
}

// A simple graph is in opposition to a multigraph: it disallows loops and
// parallel edges.
type SimpleGraph interface {
	Graph
	Density() float64
}

type DirectedGraph interface {
	Graph
	Transpose() DirectedGraph
	IsAcyclic() bool
	GetCycles() [][]Vertex
}

/* Initialization for immutable graphs */

type ImmutableGraph interface {
	Graph
	ensureVertex(vertex ...Vertex)
	addEdge(edge Edge) bool
}

var immutableGraphs map[string]ImmutableGraph

func CreateImmutableGraph(name string) (*ImmutableGraphInitializer, error) {
	template := immutableGraphs[name]
	if template == nil {
		return nil, fmt.Errorf("gogl: Unregistered graph type %s", name)
	}

	// Use reflection to make a copy of the graph template
	v := reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()
	graph, ok := v.(ImmutableGraph)
	if !ok {
		panic(fmt.Sprintf("gogl: Unable to copy graph template: %s (%v)", name, reflect.ValueOf(v).Kind().String()))
	}

	initializer := &ImmutableGraphInitializer{
		g: graph,
	}

	return initializer, nil
}

// A ImmutableGraphInitializer provides write-only methods to populate an
// immutable graph.
type ImmutableGraphInitializer struct {
	g ImmutableGraph
}

func (gi *ImmutableGraphInitializer) EnsureVertex(vertices ...Vertex) {
	gi.g.ensureVertex(vertices...)
}

func (gi *ImmutableGraphInitializer) AddEdge(edge Edge) {
	gi.g.addEdge(edge)
}

func (gi *ImmutableGraphInitializer) GetGraph() Graph {
	defer func() { gi.g = nil }()
	return gi.g
}

func init() {
	immutableGraphs = map[string]ImmutableGraph{}
}
