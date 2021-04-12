package graph

// HasNode Does the node exist in the graph
func (g *Graph) HasNode(node Node) bool {

	for index := range g.Nodes {
		if g.Nodes[index].Name == node.Name {
			return true
		}
	}

	return false
}

// UpdateNode Update the node count
func (g *Graph) UpsertNode(node *Node) {

	for index := range g.Nodes {
		if g.Nodes[index].Name == node.Name {
			g.Nodes[index].Count++
			node.Count = g.Nodes[index].Count
			node.Action = "update"
			return
		}
	}

	node.Action = "add"
	g.Nodes = append(g.Nodes, *node)

}

func (g *Graph) UpsertNodes(nodes ...Node) {
	for key := range nodes {
		g.UpsertNode(&nodes[key])
	}
}
