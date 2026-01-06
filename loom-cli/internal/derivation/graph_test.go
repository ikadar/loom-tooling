package derivation

import (
	"reflect"
	"sort"
	"testing"
)

func TestDependencyGraph_AddEdge(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-001", "BR-002", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)

	if len(g.Edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(g.Edges))
	}

	// Adding same edge again should not duplicate
	g.AddEdge("US-001", "BR-001", EdgeDerives)
	if len(g.Edges) != 3 {
		t.Errorf("Expected 3 edges (no duplicate), got %d", len(g.Edges))
	}
}

func TestDependencyGraph_RemoveEdge(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-001", "BR-002", EdgeDerives)

	g.RemoveEdge("US-001", "BR-001")

	if len(g.Edges) != 1 {
		t.Errorf("Expected 1 edge after removal, got %d", len(g.Edges))
	}

	if g.HasEdge("US-001", "BR-001") {
		t.Error("Edge should be removed")
	}

	if !g.HasEdge("US-001", "BR-002") {
		t.Error("Other edge should still exist")
	}
}

func TestDependencyGraph_RemoveNode(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)
	g.AddEdge("BR-002", "TS-001", EdgeDerives)

	g.RemoveNode("BR-001")

	if g.HasEdge("US-001", "BR-001") {
		t.Error("Incoming edge to removed node should be removed")
	}

	if g.HasEdge("BR-001", "TS-001") {
		t.Error("Outgoing edge from removed node should be removed")
	}

	if !g.HasEdge("BR-002", "TS-001") {
		t.Error("Unrelated edge should still exist")
	}
}

func TestDependencyGraph_HasEdge(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)

	if !g.HasEdge("US-001", "BR-001") {
		t.Error("Expected HasEdge to return true")
	}

	if g.HasEdge("US-001", "BR-002") {
		t.Error("Expected HasEdge to return false for non-existent edge")
	}
}

func TestDependencyGraph_GetDownstream(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-001", "BR-002", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)

	downstream := g.GetDownstream("US-001")
	sort.Strings(downstream)

	expected := []string{"BR-001", "BR-002"}
	if !reflect.DeepEqual(downstream, expected) {
		t.Errorf("Expected downstream %v, got %v", expected, downstream)
	}

	// Node with no downstream
	downstream = g.GetDownstream("TS-001")
	if len(downstream) != 0 {
		t.Error("Expected no downstream for leaf node")
	}
}

func TestDependencyGraph_GetUpstream(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-002", "BR-001", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)

	upstream := g.GetUpstream("BR-001")
	sort.Strings(upstream)

	expected := []string{"US-001", "US-002"}
	if !reflect.DeepEqual(upstream, expected) {
		t.Errorf("Expected upstream %v, got %v", expected, upstream)
	}

	// Node with no upstream
	upstream = g.GetUpstream("US-001")
	if len(upstream) != 0 {
		t.Error("Expected no upstream for root node")
	}
}

func TestDependencyGraph_GetAllDownstream(t *testing.T) {
	g := NewDependencyGraph()

	// US-001 -> BR-001 -> TS-001 -> TC-001
	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)
	g.AddEdge("TS-001", "TC-001", EdgeDerives)

	downstream := g.GetAllDownstream("US-001")
	sort.Strings(downstream)

	expected := []string{"BR-001", "TC-001", "TS-001"}
	if !reflect.DeepEqual(downstream, expected) {
		t.Errorf("Expected all downstream %v, got %v", expected, downstream)
	}
}

func TestDependencyGraph_GetAllUpstream(t *testing.T) {
	g := NewDependencyGraph()

	// US-001 -> BR-001 -> TS-001 -> TC-001
	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)
	g.AddEdge("TS-001", "TC-001", EdgeDerives)

	upstream := g.GetAllUpstream("TC-001")
	sort.Strings(upstream)

	expected := []string{"BR-001", "TS-001", "US-001"}
	if !reflect.DeepEqual(upstream, expected) {
		t.Errorf("Expected all upstream %v, got %v", expected, upstream)
	}
}

func TestDependencyGraph_GetRoots(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-002", "BR-002", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)

	roots := g.GetRoots()

	expected := []string{"US-001", "US-002"}
	if !reflect.DeepEqual(roots, expected) {
		t.Errorf("Expected roots %v, got %v", expected, roots)
	}
}

func TestDependencyGraph_GetLeaves(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)
	g.AddEdge("BR-001", "TS-002", EdgeDerives)

	leaves := g.GetLeaves()

	expected := []string{"TS-001", "TS-002"}
	if !reflect.DeepEqual(leaves, expected) {
		t.Errorf("Expected leaves %v, got %v", expected, leaves)
	}
}

func TestDependencyGraph_TopologicalSort(t *testing.T) {
	g := NewDependencyGraph()

	// Build a DAG
	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-001", "BR-002", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)
	g.AddEdge("BR-002", "TS-001", EdgeDerives)
	g.AddEdge("TS-001", "TC-001", EdgeDerives)

	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify order constraints
	indexOf := func(id string) int {
		for i, v := range order {
			if v == id {
				return i
			}
		}
		return -1
	}

	// US-001 must come before BR-001 and BR-002
	if indexOf("US-001") > indexOf("BR-001") || indexOf("US-001") > indexOf("BR-002") {
		t.Error("US-001 should come before BR-001 and BR-002")
	}

	// BR-001 and BR-002 must come before TS-001
	if indexOf("BR-001") > indexOf("TS-001") || indexOf("BR-002") > indexOf("TS-001") {
		t.Error("BR-001 and BR-002 should come before TS-001")
	}

	// TS-001 must come before TC-001
	if indexOf("TS-001") > indexOf("TC-001") {
		t.Error("TS-001 should come before TC-001")
	}
}

func TestDependencyGraph_TopologicalSort_Cycle(t *testing.T) {
	g := NewDependencyGraph()

	// Create a cycle: A -> B -> C -> A
	g.AddEdge("A", "B", EdgeDerives)
	g.AddEdge("B", "C", EdgeDerives)
	g.AddEdge("C", "A", EdgeDerives)

	_, err := g.TopologicalSort()
	if err == nil {
		t.Error("Expected error for cyclic graph")
	}
}

func TestDependencyGraph_DetectCycle(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		g := NewDependencyGraph()
		g.AddEdge("A", "B", EdgeDerives)
		g.AddEdge("B", "C", EdgeDerives)

		cycle := g.DetectCycle()
		if cycle != nil {
			t.Errorf("Expected no cycle, got %v", cycle)
		}
	})

	t.Run("simple cycle", func(t *testing.T) {
		g := NewDependencyGraph()
		g.AddEdge("A", "B", EdgeDerives)
		g.AddEdge("B", "C", EdgeDerives)
		g.AddEdge("C", "A", EdgeDerives)

		cycle := g.DetectCycle()
		if cycle == nil {
			t.Error("Expected to detect cycle")
		}
	})

	t.Run("self loop", func(t *testing.T) {
		g := NewDependencyGraph()
		g.AddEdge("A", "A", EdgeDerives)

		cycle := g.DetectCycle()
		if cycle == nil {
			t.Error("Expected to detect self-loop")
		}
	})
}

func TestDependencyGraph_RebuildFromEdges(t *testing.T) {
	g := NewDependencyGraph()

	// Add edges directly to Edges slice (simulating JSON load)
	g.Edges = []DependencyEdge{
		{From: "A", To: "B", Type: EdgeDerives},
		{From: "B", To: "C", Type: EdgeDerives},
	}

	// Before rebuild, adjacency should be empty
	if len(g.adjacency) != 0 {
		t.Error("Expected empty adjacency before rebuild")
	}

	g.RebuildFromEdges()

	// After rebuild, should be able to query
	downstream := g.GetDownstream("A")
	if len(downstream) != 1 || downstream[0] != "B" {
		t.Error("Expected downstream B from A after rebuild")
	}
}

func TestDependencyGraph_BuildFromArtifacts(t *testing.T) {
	g := NewDependencyGraph()

	artifacts := map[string]*Artifact{
		"BR-001": {
			ID: "BR-001",
			Upstream: map[string]string{
				"US-001": "hash1",
			},
		},
		"TS-001": {
			ID: "TS-001",
			Upstream: map[string]string{
				"BR-001": "hash2",
			},
		},
	}

	g.BuildFromArtifacts(artifacts)

	if !g.HasEdge("US-001", "BR-001") {
		t.Error("Expected edge US-001 -> BR-001")
	}
	if !g.HasEdge("BR-001", "TS-001") {
		t.Error("Expected edge BR-001 -> TS-001")
	}
}

func TestDependencyGraph_NodeCount(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("A", "B", EdgeDerives)
	g.AddEdge("B", "C", EdgeDerives)
	g.AddEdge("A", "C", EdgeDerives)

	if g.NodeCount() != 3 {
		t.Errorf("Expected 3 nodes, got %d", g.NodeCount())
	}
}

func TestDependencyGraph_EdgeCount(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("A", "B", EdgeDerives)
	g.AddEdge("B", "C", EdgeDerives)
	g.AddEdge("A", "C", EdgeDerives)

	if g.EdgeCount() != 3 {
		t.Errorf("Expected 3 edges, got %d", g.EdgeCount())
	}
}

func TestDependencyGraph_GetNodes(t *testing.T) {
	g := NewDependencyGraph()

	g.AddEdge("C", "D", EdgeDerives)
	g.AddEdge("A", "B", EdgeDerives)

	nodes := g.GetNodes()
	expected := []string{"A", "B", "C", "D"}

	if !reflect.DeepEqual(nodes, expected) {
		t.Errorf("Expected nodes %v, got %v", expected, nodes)
	}
}

func TestDependencyGraph_GetAffectedByChange(t *testing.T) {
	g := NewDependencyGraph()

	// US-001 -> BR-001 -> TS-001
	//        -> BR-002 -> TS-002
	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-001", "BR-002", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)
	g.AddEdge("BR-002", "TS-002", EdgeDerives)

	affected := g.GetAffectedByChange("US-001")
	sort.Strings(affected)

	expected := []string{"BR-001", "BR-002", "TS-001", "TS-002"}
	if !reflect.DeepEqual(affected, expected) {
		t.Errorf("Expected affected %v, got %v", expected, affected)
	}
}

func TestDependencyGraph_GetDerivationOrder(t *testing.T) {
	g := NewDependencyGraph()

	// US-001 -> BR-001 -> TS-001 -> TC-001
	//        -> BR-002 -> TS-002
	g.AddEdge("US-001", "BR-001", EdgeDerives)
	g.AddEdge("US-001", "BR-002", EdgeDerives)
	g.AddEdge("BR-001", "TS-001", EdgeDerives)
	g.AddEdge("BR-002", "TS-002", EdgeDerives)
	g.AddEdge("TS-001", "TC-001", EdgeDerives)

	// If BR-001 is stale, derivation order should include BR-001, TS-001, TC-001
	order, err := g.GetDerivationOrder([]string{"BR-001"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify BR-001 is in the list
	found := false
	for _, id := range order {
		if id == "BR-001" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected BR-001 in derivation order")
	}

	// Verify TS-001 and TC-001 are included (affected)
	expectedIDs := map[string]bool{"BR-001": true, "TS-001": true, "TC-001": true}
	for _, id := range order {
		if !expectedIDs[id] {
			t.Errorf("Unexpected ID %s in derivation order", id)
		}
	}

	if len(order) != 3 {
		t.Errorf("Expected 3 artifacts in order, got %d", len(order))
	}
}

func TestDependencyGraph_Empty(t *testing.T) {
	g := NewDependencyGraph()

	if g.NodeCount() != 0 {
		t.Error("Expected 0 nodes in empty graph")
	}
	if g.EdgeCount() != 0 {
		t.Error("Expected 0 edges in empty graph")
	}

	roots := g.GetRoots()
	if len(roots) != 0 {
		t.Error("Expected no roots in empty graph")
	}

	order, err := g.TopologicalSort()
	if err != nil {
		t.Error("TopologicalSort should succeed on empty graph")
	}
	if len(order) != 0 {
		t.Error("Expected empty order for empty graph")
	}
}
