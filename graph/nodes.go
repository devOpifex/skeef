package graph

// HasNode Does the node exist in the graph
func (g *Graph) HasNode(name string) bool {

	for index := range g.Nodes {
		if g.Nodes[index].Name == name {
			return true
		}
	}

	return false
}

// UpdateNode Update the node count
func (g *Graph) UpdateNode(name string, count int) {

	for index := range g.Nodes {
		if g.Nodes[index].Name == name {
			g.Nodes[index].Count = count
			break
		}
	}

}
