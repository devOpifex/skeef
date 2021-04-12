package graph

import (
	"github.com/dghubble/go-twitter/twitter"
)

// Node defines nodes
type Node struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Count  int    `json:"count"`
	Action string `json:"action"`
}

// Edge edges
type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Weight int    `json:"weight"`
	Action string `json:"action"`
}

// Graph defines a graph
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// GetUserNet builds the network of users, where one user
// mentions another
func GetUserNet(tweet twitter.Tweet, exclusion map[string]bool) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, m := range tweet.Entities.UserMentions {
		_, ok := exclusion[tweet.User.ScreenName]

		if ok {
			continue
		}

		_, ok = exclusion[m.ScreenName]

		if ok {
			continue
		}

		edge := Edge{tweet.User.ScreenName, m.ScreenName, 1, "add"}

		edges = append(edges, edge)
	}

	for _, n := range edges {
		src := Node{n.Source, "user", 1, "add"}
		tgt := Node{n.Target, "user", 1, "add"}
		nodes = append(nodes, src, tgt)
	}

	return nodes, edges
}

// GetHashNet builds the network of users to hashtags they use in tweets
func GetHashNet(tweet twitter.Tweet, exclusion map[string]bool) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, h := range tweet.Entities.Hashtags {

		_, ok := exclusion[tweet.User.ScreenName]

		if ok {
			continue
		}

		_, ok = exclusion[h.Text]

		if ok {
			continue
		}

		edge := Edge{tweet.User.ScreenName, h.Text, 1, "add"}
		edges = append(edges, edge)
	}

	for _, e := range edges {
		src := Node{e.Source, "user", 1, "add"}
		tgt := Node{e.Target, "hashtag", 1, "add"}

		nodes = append(nodes, src, tgt)
	}

	return nodes, edges
}

// GetUserNet builds the network of users, where one user
// mentions another
func GetRetweetNet(tweet twitter.Tweet, exclusion map[string]bool) (bool, []Node, Edge) {

	var edge Edge
	var nodes []Node

	if tweet.InReplyToScreenName == "" {
		return false, nodes, edge
	}

	_, ok := exclusion[tweet.InReplyToScreenName]

	if ok {
		return false, nodes, edge
	}

	_, ok = exclusion[tweet.User.ScreenName]

	if ok {
		return false, nodes, edge
	}

	edge = Edge{tweet.User.ScreenName, tweet.InReplyToScreenName, 1, "add"}
	from := Node{tweet.User.ScreenName, "user", 1, "add"}
	to := Node{tweet.InReplyToScreenName, "user", 1, "add"}

	nodes = append(nodes, from, to)

	return true, nodes, edge
}

func (g *Graph) Truncate(max int) ([]Node, []Edge) {

	var edgesToRemove []Edge
	var nodesToRemove []Node

	// no need to truncate
	if len(g.Edges) < max {
		return nodesToRemove, edgesToRemove
	}

	diff := len(g.Edges) - max
	edgesToRemove = g.Edges[0:diff]
	for i := range edgesToRemove {
		edgesToRemove[i].Action = "remove"
	}

	var nodesKeep []Node
	for _, n := range g.Nodes {
		for _, e := range g.Edges {

			if n.Name == e.Source || n.Name == e.Target {
				n.Action = "remove"
				nodesToRemove = append(nodesToRemove, n)
			} else {
				nodesKeep = append(nodesKeep, n)
			}

		}
	}

	g.Nodes = nodesKeep

	return nodesToRemove, edgesToRemove
}
