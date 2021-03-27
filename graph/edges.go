package graph

// HasNode Does the node exist in the graph
func (g *Graph) HasEdge(source, target string) bool {

	for _, edge := range g.Edges {
		if edge.Source == source && edge.Target == target {
			return true
		}
	}

	return false
}

// UpdateNode Update the node count
func (g *Graph) UpdateEdge(source, target string, weight int) {

	for index := range g.Edges {
		if g.Edges[index].Source == source && g.Edges[index].Target == target {
			g.Edges[index].Weight = weight
			break
		}
	}

}
