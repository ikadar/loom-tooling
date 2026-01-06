package derivation

import (
	"fmt"
	"sort"
	"sync"
)

// =============================================================================
// Dependency Graph
// =============================================================================

// DependencyGraph represents the directed acyclic graph of artifact dependencies
type DependencyGraph struct {
	// Edges stores all dependency relationships
	Edges []DependencyEdge `json:"edges"`

	// adjacency maps artifact ID to its downstream dependencies
	adjacency map[string][]string

	// reverse maps artifact ID to its upstream dependencies
	reverse map[string][]string

	// mu protects concurrent access
	mu sync.RWMutex
}

// NewDependencyGraph creates a new empty dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Edges:     make([]DependencyEdge, 0),
		adjacency: make(map[string][]string),
		reverse:   make(map[string][]string),
	}
}

// =============================================================================
// Edge Operations
// =============================================================================

// AddEdge adds a dependency edge from upstream to downstream
func (g *DependencyGraph) AddEdge(from, to string, edgeType EdgeType) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check if edge already exists
	for _, e := range g.Edges {
		if e.From == from && e.To == to {
			return // Edge already exists
		}
	}

	// Add edge
	g.Edges = append(g.Edges, DependencyEdge{
		From: from,
		To:   to,
		Type: edgeType,
	})

	// Update adjacency lists
	g.adjacency[from] = append(g.adjacency[from], to)
	g.reverse[to] = append(g.reverse[to], from)
}

// RemoveEdge removes a dependency edge
func (g *DependencyGraph) RemoveEdge(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Remove from edges list
	for i, e := range g.Edges {
		if e.From == from && e.To == to {
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
			break
		}
	}

	// Update adjacency lists
	g.adjacency[from] = removeFromSlice(g.adjacency[from], to)
	g.reverse[to] = removeFromSlice(g.reverse[to], from)
}

// RemoveNode removes a node and all its edges
func (g *DependencyGraph) RemoveNode(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Remove all edges involving this node
	newEdges := make([]DependencyEdge, 0, len(g.Edges))
	for _, e := range g.Edges {
		if e.From != id && e.To != id {
			newEdges = append(newEdges, e)
		}
	}
	g.Edges = newEdges

	// Update adjacency lists
	// Remove from downstream lists
	for _, downstream := range g.adjacency[id] {
		g.reverse[downstream] = removeFromSlice(g.reverse[downstream], id)
	}
	delete(g.adjacency, id)

	// Remove from upstream lists
	for _, upstream := range g.reverse[id] {
		g.adjacency[upstream] = removeFromSlice(g.adjacency[upstream], id)
	}
	delete(g.reverse, id)
}

// HasEdge checks if an edge exists
func (g *DependencyGraph) HasEdge(from, to string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, e := range g.Edges {
		if e.From == from && e.To == to {
			return true
		}
	}
	return false
}

// =============================================================================
// Query Operations
// =============================================================================

// GetDownstream returns direct downstream dependencies of an artifact
func (g *DependencyGraph) GetDownstream(id string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make([]string, len(g.adjacency[id]))
	copy(result, g.adjacency[id])
	return result
}

// GetUpstream returns direct upstream dependencies of an artifact
func (g *DependencyGraph) GetUpstream(id string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make([]string, len(g.reverse[id]))
	copy(result, g.reverse[id])
	return result
}

// GetAllDownstream returns all downstream dependencies (transitive closure)
func (g *DependencyGraph) GetAllDownstream(id string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	var result []string

	var dfs func(current string)
	dfs = func(current string) {
		for _, next := range g.adjacency[current] {
			if !visited[next] {
				visited[next] = true
				result = append(result, next)
				dfs(next)
			}
		}
	}

	dfs(id)
	return result
}

// GetAllUpstream returns all upstream dependencies (transitive closure)
func (g *DependencyGraph) GetAllUpstream(id string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	var result []string

	var dfs func(current string)
	dfs = func(current string) {
		for _, prev := range g.reverse[current] {
			if !visited[prev] {
				visited[prev] = true
				result = append(result, prev)
				dfs(prev)
			}
		}
	}

	dfs(id)
	return result
}

// GetRoots returns all nodes with no upstream dependencies
func (g *DependencyGraph) GetRoots() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Find all nodes
	nodes := make(map[string]bool)
	for _, e := range g.Edges {
		nodes[e.From] = true
		nodes[e.To] = true
	}

	// Find roots (no incoming edges)
	var roots []string
	for node := range nodes {
		if len(g.reverse[node]) == 0 {
			roots = append(roots, node)
		}
	}

	sort.Strings(roots)
	return roots
}

// GetLeaves returns all nodes with no downstream dependencies
func (g *DependencyGraph) GetLeaves() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Find all nodes
	nodes := make(map[string]bool)
	for _, e := range g.Edges {
		nodes[e.From] = true
		nodes[e.To] = true
	}

	// Find leaves (no outgoing edges)
	var leaves []string
	for node := range nodes {
		if len(g.adjacency[node]) == 0 {
			leaves = append(leaves, node)
		}
	}

	sort.Strings(leaves)
	return leaves
}

// =============================================================================
// Topological Sort
// =============================================================================

// TopologicalSort returns artifacts in dependency order (upstream before downstream)
// Returns error if the graph contains a cycle
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Find all nodes
	nodes := make(map[string]bool)
	for _, e := range g.Edges {
		nodes[e.From] = true
		nodes[e.To] = true
	}

	// Calculate in-degrees
	inDegree := make(map[string]int)
	for node := range nodes {
		inDegree[node] = len(g.reverse[node])
	}

	// Start with nodes that have no dependencies
	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}
	sort.Strings(queue) // Ensure deterministic order

	var result []string
	for len(queue) > 0 {
		// Pop first element
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// Reduce in-degree for downstream nodes
		for _, downstream := range g.adjacency[node] {
			inDegree[downstream]--
			if inDegree[downstream] == 0 {
				queue = append(queue, downstream)
				sort.Strings(queue) // Maintain deterministic order
			}
		}
	}

	// Check for cycle
	if len(result) != len(nodes) {
		return nil, fmt.Errorf("dependency graph contains a cycle")
	}

	return result, nil
}

// =============================================================================
// Cycle Detection
// =============================================================================

// DetectCycle checks if the graph contains a cycle
// Returns the cycle path if found, nil otherwise
func (g *DependencyGraph) DetectCycle() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Find all nodes
	nodes := make(map[string]bool)
	for _, e := range g.Edges {
		nodes[e.From] = true
		nodes[e.To] = true
	}

	// DFS with coloring
	// 0 = white (unvisited), 1 = gray (visiting), 2 = black (visited)
	color := make(map[string]int)
	parent := make(map[string]string)

	var cyclePath []string

	var dfs func(node string) bool
	dfs = func(node string) bool {
		color[node] = 1 // Gray - visiting

		for _, next := range g.adjacency[node] {
			if color[next] == 1 {
				// Found cycle - reconstruct path
				cyclePath = []string{next, node}
				current := node
				for parent[current] != "" && parent[current] != next {
					current = parent[current]
					cyclePath = append([]string{current}, cyclePath...)
				}
				return true
			}

			if color[next] == 0 {
				parent[next] = node
				if dfs(next) {
					return true
				}
			}
		}

		color[node] = 2 // Black - visited
		return false
	}

	for node := range nodes {
		if color[node] == 0 {
			if dfs(node) {
				return cyclePath
			}
		}
	}

	return nil
}

// =============================================================================
// Graph Building
// =============================================================================

// RebuildFromEdges rebuilds the adjacency lists from the edges
// Call this after loading state from JSON
func (g *DependencyGraph) RebuildFromEdges() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.adjacency = make(map[string][]string)
	g.reverse = make(map[string][]string)

	for _, e := range g.Edges {
		g.adjacency[e.From] = append(g.adjacency[e.From], e.To)
		g.reverse[e.To] = append(g.reverse[e.To], e.From)
	}
}

// BuildFromArtifacts builds the graph from artifact upstream/downstream data
func (g *DependencyGraph) BuildFromArtifacts(artifacts map[string]*Artifact) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Edges = make([]DependencyEdge, 0)
	g.adjacency = make(map[string][]string)
	g.reverse = make(map[string][]string)

	for _, artifact := range artifacts {
		// Add edges from upstream dependencies
		for upstreamID := range artifact.Upstream {
			g.Edges = append(g.Edges, DependencyEdge{
				From: upstreamID,
				To:   artifact.ID,
				Type: EdgeDerives,
			})
			g.adjacency[upstreamID] = append(g.adjacency[upstreamID], artifact.ID)
			g.reverse[artifact.ID] = append(g.reverse[artifact.ID], upstreamID)
		}
	}
}

// =============================================================================
// Graph Statistics
// =============================================================================

// NodeCount returns the number of nodes in the graph
func (g *DependencyGraph) NodeCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	nodes := make(map[string]bool)
	for _, e := range g.Edges {
		nodes[e.From] = true
		nodes[e.To] = true
	}
	return len(nodes)
}

// EdgeCount returns the number of edges in the graph
func (g *DependencyGraph) EdgeCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.Edges)
}

// GetNodes returns all node IDs in the graph
func (g *DependencyGraph) GetNodes() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	nodes := make(map[string]bool)
	for _, e := range g.Edges {
		nodes[e.From] = true
		nodes[e.To] = true
	}

	result := make([]string, 0, len(nodes))
	for node := range nodes {
		result = append(result, node)
	}
	sort.Strings(result)
	return result
}

// =============================================================================
// Impact Analysis
// =============================================================================

// GetAffectedByChange returns all artifacts that would be affected if the given artifact changes
// This includes direct downstream and transitive downstream
func (g *DependencyGraph) GetAffectedByChange(changedID string) []string {
	return g.GetAllDownstream(changedID)
}

// GetDerivationOrder returns the order in which artifacts should be derived
// given a set of stale artifact IDs
func (g *DependencyGraph) GetDerivationOrder(staleIDs []string) ([]string, error) {
	// Get all affected artifacts (including stale)
	affected := make(map[string]bool)
	for _, id := range staleIDs {
		affected[id] = true
		for _, downstream := range g.GetAllDownstream(id) {
			affected[downstream] = true
		}
	}

	// Get topological order of all nodes
	allOrder, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}

	// Filter to only affected artifacts
	var result []string
	for _, id := range allOrder {
		if affected[id] {
			result = append(result, id)
		}
	}

	return result, nil
}

// =============================================================================
// Helper Functions
// =============================================================================

// removeFromSlice removes the first occurrence of a value from a slice
func removeFromSlice(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
