package graph

import (
	"github.com/dghubble/go-twitter/twitter"
)

// GetUserNet builds the network of users, where one user
// mentions another
func GetUserNet(tweet twitter.Tweet) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, m := range tweet.Entities.UserMentions {
		edge := Edge{tweet.User.ScreenName, m.ScreenName, 0}

		edges = append(edges, edge)
	}

	for _, n := range edges {
		src := Node{n.Source, "user", 0}
		tgt := Node{n.Target, "user", 0}
		nodes = append(nodes, src, tgt)
	}

	return nodes, edges
}

// GetHashNet builds the network of users to hashtags they use in tweets
func GetHashNet(tweet twitter.Tweet) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, h := range tweet.Entities.Hashtags {

		edge := Edge{tweet.User.ScreenName, h.Text, 0}
		edges = append(edges, edge)
	}

	for _, n := range edges {
		src := Node{n.Source, "user", 0}
		tgt := Node{n.Target, "hashtag", 0}

		nodes = append(nodes, src, tgt)
	}

	return nodes, edges
}

// BuildGraph builds the graph from edges and nodes
func BuildGraph(nodes []Node, edges []Edge) Graph {

	var g Graph
	var newNodes []Node
	var users = make(map[string]int)
	var hashtags = make(map[string]int)

	// handle nodes
	for _, n := range nodes {
		if n.Type == "hashtag" {
			hashtags[n.Name]++
		}

		if n.Type == "user" {
			users[n.Name]++
		}
	}

	for name, count := range users {
		n := Node{name, "user", count}
		newNodes = append(newNodes, n)
	}

	for hash, count := range hashtags {
		n := Node{hash, "hashtag", count}
		newNodes = append(newNodes, n)
	}

	g.Nodes = newNodes
	g.Edges = edges

	return g
}
