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
func GetMentionNet(tweet twitter.Tweet, exclusion map[string]bool, minFollowerCount, minFavoriteCount int, onlyVerified bool, maxHashtags, maxMentions int) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, m := range tweet.Entities.UserMentions {

		// filters
		_, ok := exclusion[tweet.User.ScreenName]

		if ok {
			continue
		}

		if tweet.User.FollowersCount < minFollowerCount {
			continue
		}

		if tweet.FavoriteCount < minFavoriteCount {
			continue
		}

		if onlyVerified && !tweet.User.Verified {
			continue
		}

		if len(tweet.Entities.Hashtags) > maxHashtags {
			continue
		}

		if len(tweet.Entities.UserMentions) > maxMentions {
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
func GetHashtagNet(tweet twitter.Tweet, exclusion map[string]bool, minFollowerCount, minFavoriteCount int, onlyVerified bool, maxHashtags, maxMentions int) ([]Node, []Edge) {

	var edges []Edge
	var nodes []Node

	for _, h := range tweet.Entities.Hashtags {

		_, ok := exclusion[tweet.User.ScreenName]

		if ok {
			continue
		}

		_, ok = exclusion["#"+h.Text]

		if ok {
			continue
		}

		if tweet.User.FollowersCount <= minFollowerCount {
			continue
		}

		if tweet.FavoriteCount <= minFavoriteCount {
			continue
		}

		if onlyVerified && !tweet.User.Verified {
			continue
		}

		if len(tweet.Entities.Hashtags) >= maxHashtags {
			continue
		}

		if len(tweet.Entities.UserMentions) >= maxMentions {
			continue
		}

		edge := Edge{tweet.User.ScreenName, "#" + h.Text, 1, "add"}
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
func GetRetweetNet(tweet twitter.Tweet, exclusion map[string]bool, minFollowerCount, minFavoriteCount int, onlyVerified bool, maxHashtags, maxMentions int) (bool, []Node, Edge) {

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

	if tweet.User.FollowersCount <= minFollowerCount {
		return false, nodes, edge
	}

	if tweet.FavoriteCount <= minFavoriteCount {
		return false, nodes, edge
	}

	if onlyVerified && !tweet.User.Verified {
		return false, nodes, edge
	}

	if len(tweet.Entities.Hashtags) >= maxHashtags {
		return false, nodes, edge
	}

	if len(tweet.Entities.UserMentions) >= maxMentions {
		return false, nodes, edge
	}

	edge = Edge{tweet.User.ScreenName, tweet.InReplyToScreenName, 1, "add"}
	from := Node{tweet.User.ScreenName, "user", 1, "add"}
	to := Node{tweet.InReplyToScreenName, "user", 1, "add"}

	nodes = append(nodes, from, to)

	return true, nodes, edge
}

// Truncate truncates the graph to ensure we limit the number of
// edges (and subsequently nodes) present on the screen
// at any one time.
func (g *Graph) Truncate(max int) ([]Node, []Edge) {

	var nodesToRemove []Node
	var edgesToRemove []Edge

	// no need to truncate we're below the threshold
	if len(g.Edges) < max {
		return nodesToRemove, edgesToRemove
	}

	// keep track of edges that we should remove
	// a map is much more efficient
	edgesToRemoveMap := make(map[string]bool)

	diff := len(g.Edges) - max
	edgesToRemove = g.Edges[0:diff]

	// no edges to remove
	if len(edgesToRemove) < 1 {
		return nodesToRemove, edgesToRemove
	}

	// convert the struct of edges into a map
	for _, e := range edgesToRemove {
		edgesToRemoveMap[e.Source] = true
		edgesToRemoveMap[e.Target] = true
	}

	// we only keep those edges
	g.Edges = g.Edges[diff:len(g.Edges)]

	// we cannot delete a node that is part
	// of an edge that remains on the graph
	for _, e := range g.Edges {
		edgesToRemoveMap[e.Source] = false
		edgesToRemoveMap[e.Target] = false
	}

	var nodesKeep []Node
	for _, n := range g.Nodes {
		ok := edgesToRemoveMap[n.Name]

		if ok {
			nodesToRemove = append(nodesToRemove, n)
		} else {
			nodesKeep = append(nodesKeep, n)
		}
	}

	g.Nodes = nodesKeep

	return nodesToRemove, edgesToRemove
}

// GetReplyNet builds an edge that connects the author of a tweet
// with the user the author responds to
func GetReplyNet(tweet twitter.Tweet, exclusion map[string]bool, minFollowerCount, minFavoriteCount int, onlyVerified bool, maxHashtags, maxMentions int) ([]Node, Edge) {

	var edge Edge
	var nodes []Node

	// filters
	_, ok := exclusion[tweet.InReplyToScreenName]

	if ok {
		return nodes, edge
	}

	if tweet.User.FollowersCount < minFollowerCount {
		return nodes, edge
	}

	if tweet.FavoriteCount < minFavoriteCount {
		return nodes, edge
	}

	if onlyVerified && !tweet.User.Verified {
		return nodes, edge
	}

	if len(tweet.Entities.Hashtags) > maxHashtags {
		return nodes, edge
	}

	if len(tweet.Entities.UserMentions) > maxMentions {
		return nodes, edge
	}

	edge = Edge{tweet.User.ScreenName, tweet.InReplyToScreenName, 1, "add"}

	src := Node{edge.Source, "user", 1, "add"}
	tgt := Node{edge.Target, "user", 1, "add"}
	nodes = append(nodes, src, tgt)

	return nodes, edge
}
