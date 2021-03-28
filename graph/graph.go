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
func GetUserNet(tweet twitter.Tweet) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, m := range tweet.Entities.UserMentions {
		edge := Edge{tweet.User.ScreenName, m.ScreenName, 1, ""}

		edges = append(edges, edge)
	}

	for _, n := range edges {
		src := Node{n.Source, "user", 1, ""}
		tgt := Node{n.Target, "user", 1, ""}
		nodes = append(nodes, src, tgt)
	}

	return nodes, edges
}

// GetHashNet builds the network of users to hashtags they use in tweets
func GetHashNet(tweet twitter.Tweet) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, h := range tweet.Entities.Hashtags {

		edge := Edge{tweet.User.ScreenName, h.Text, 1, ""}
		edges = append(edges, edge)
	}

	for _, e := range edges {
		src := Node{e.Source, "user", 1, ""}
		tgt := Node{e.Target, "hashtag", 1, ""}

		nodes = append(nodes, src, tgt)
	}

	return nodes, edges
}
