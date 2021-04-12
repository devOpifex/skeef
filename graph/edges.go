package graph

// HasNode Does the node exist in the graph
func (g *Graph) HasEdge(edge Edge) bool {

	for _, e := range g.Edges {
		if e.Source == edge.Source && e.Target == edge.Target {
			return true
		}
	}

	return false
}

// UpdateNode Update the node count
func (g *Graph) UpsertEdge(edge *Edge) {

	for index := range g.Edges {
		if g.Edges[index].Source == edge.Source && g.Edges[index].Target == edge.Target {
			g.Edges[index].Weight++
			edge.Weight = g.Edges[index].Weight
			edge.Action = "update"
			return
		}
	}

	edge.Action = "add"
	g.Edges = append(g.Edges, *edge)

}

func (g *Graph) UpsertEdges(edges ...Edge) {
	for key := range edges {
		g.UpsertEdge(&edges[key])
	}
}
