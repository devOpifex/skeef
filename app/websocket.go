package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/devOpifex/skeef-app/graph"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/websocket"
)

type message struct {
	Graph       graph.Graph   `json:"graph"`
	Trend       map[int64]int `json:"trend"`
	RemoveEdges []graph.Edge  `json:"removeEdges"`
	RemoveNodes []graph.Node  `json:"removeNodes"`
	NewTweets   tweetsUsers   `json:"tweets"`
}

type connectionMessage struct {
	Graph graph.Graph   `json:"graph"`
	Trend map[int64]int `json:"trend"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan message),
	}
}

const maxMessageSize = 512

func (c *Client) Read(app *Application) {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

	}
}

func (app *Application) StartPool() {
	for {
		select {
		case client := <-app.Pool.Register:
			app.Pool.Clients[client] = true
			app.Connected++
			for client := range app.Pool.Clients {
				// send current state of the graph on connect
				client.Conn.WriteJSON(
					connectionMessage{
						Graph: app.Graph,
						Trend: app.Trend,
					})
			}
		case client := <-app.Pool.Unregister:
			delete(app.Pool.Clients, client)
			app.Connected--
			// for client := range app.Pool.Clients {
			// 	client.Conn.WriteJSON(message{})
			// }
		case message := <-app.Pool.Broadcast:
			for client := range app.Pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// wsUpgrade Upgrades the websocket
func (app *Application) wsUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.ErrorLog.Println(err)
		return ws, err
	}
	return ws, nil
}

// socket Upgrade the websocket and app it to the pool
func (app *Application) socket(w http.ResponseWriter, r *http.Request) {
	ws, err := app.wsUpgrade(w, r)

	if err != nil {
		log.Println(err)
	}

	client := &Client{
		Conn: ws,
		Pool: app.Pool,
	}

	app.Pool.Register <- client
	client.Read(app)

}

// StopStream stops the stream
func (app *Application) StopStream() {

	if app.Stream != nil {
		app.Stream.Stop()
	}

	app.Quit <- struct{}{}
	app.Database.PauseAllStreams()
}

// StartStream starts the stream
func (app *Application) StartStream() {

	app.Trend = make(map[int64]int)
	app.Count = 0
	app.Graph = graph.Graph{}

	tokens, err := app.Database.GetTokens()

	if err != nil {
		return
	}

	app.StreamActive = app.Database.GetActiveStream()
	app.Exclusion = exclusionMap(app.StreamActive.Exclusion)
	app.MaxEdges = app.StreamActive.MaxEdges

	var twitterConfig = oauth1.NewConfig(tokens.ApiKey, tokens.ApiSecret)
	var token = oauth1.NewToken(tokens.AccessToken, tokens.AccessSecret)

	var httpClient = twitterConfig.Client(oauth1.NoContext, token)

	var client = twitter.NewClient(httpClient)

	var params = &twitter.StreamFilterParams{
		Track:         splitTerm(app.StreamActive.Track),
		Follow:        splitTerm(app.StreamActive.Follow),
		StallWarnings: twitter.Bool(true),
	}
	app.Stream, _ = client.Streams.Filter(params)

	var demux = twitter.NewSwitchDemux()
	demux.Tweet = app.demux()

	for message := range app.Stream.Messages {
		demux.Handle(message)
	}
}

// demux Demux tweets and broadcast to websocket clients
func (app *Application) demux() func(tweet *twitter.Tweet) {
	return func(tweet *twitter.Tweet) {
		app.Count++
		app.Trend[parseTime(tweet.CreatedAt)]++
		app.truncateTrend()
		app.InfoLog.Printf("Tweet #%v\n", app.Count)

		// declare variables
		var nodes []graph.Node
		var edges []graph.Edge

		// selectively create graph
		if app.StreamActive.MentionsNet > 0 {
			nodesMentions, edgesMentions := graph.GetUserNet(*tweet, app.Exclusion)
			nodes = append(nodes, nodesMentions...)
			edges = append(edges, edgesMentions...)
		}
		if app.StreamActive.HashtagsNet > 0 {
			nodesHash, edgesHash := graph.GetHashNet(*tweet, app.Exclusion)
			nodes = append(nodes, nodesHash...)
			edges = append(edges, edgesHash...)
		}
		if app.StreamActive.RetweetsNet > 0 {
			ok, nodesRetweet, edgesRetweet := graph.GetRetweetNet(*tweet, app.Exclusion)
			if ok {
				nodes = append(nodes, nodesRetweet...)
				edges = append(edges, edgesRetweet)
			}
		}
		newTweets := app.trackTweets(tweet, edges)

		app.Graph.UpsertEdges(edges...)
		app.Graph.UpsertNodes(nodes...)

		removeNodes, removeEdges := app.Graph.Truncate(app.MaxEdges)

		g := graph.Graph{
			Edges: edges,
			Nodes: nodes,
		}

		m := message{
			Graph:       g,
			Trend:       app.Trend,
			RemoveNodes: removeNodes,
			RemoveEdges: removeEdges,
			NewTweets:   newTweets,
		}
		app.Pool.Broadcast <- m
	}
}

// parseTime Parse the ruby date and round to nearest 15 seconds
func parseTime(date string) int64 {
	toRound, _ := time.Parse(time.RubyDate, date)
	minute := toRound.Round(15 * time.Second)

	return minute.Unix()
}

// splitTerm splits comma separate string into
// slice of strings
func splitTerm(track string) []string {
	splat := strings.Split(track, ",")

	for i, s := range splat {
		splat[i] = strings.TrimSpace(s)
	}

	return splat
}

// truncateTrend Truncates the trend to ensure it
// does not grow infinitely
func (app *Application) truncateTrend() {
	if len(app.Trend) < 50 {
		return
	}

	var min int64

	for key := range app.Trend {
		if min == 0 {
			min = key
			continue
		}

		if min < key {
			continue
		}

		min = key
	}

	delete(app.Trend, min)
}

type tweetsUsers map[string]map[string]string

func (app *Application) trackTweets(tweet *twitter.Tweet, edges []graph.Edge) tweetsUsers {
	var newTweets = make(map[string]map[string]string)
	newTweets[tweet.IDStr] = make(map[string]string)

	for _, v := range edges {

		// need to create map
		if _, ok := app.TweetsUsers[tweet.IDStr]; !ok {
			app.TweetsUsers[tweet.IDStr] = make(map[string]string)
		}

		// newly added tweets and/or
		if _, ok := app.TweetsUsers[tweet.IDStr][v.Source]; !ok {
			newTweets[tweet.IDStr][v.Source] = v.Source
		}
		if _, ok := app.TweetsUsers[tweet.IDStr][v.Target]; !ok {
			newTweets[tweet.IDStr][v.Target] = v.Target
		}

		// global for initial connection
		app.TweetsUsers[tweet.IDStr][v.Source] = v.Source
		app.TweetsUsers[tweet.IDStr][v.Target] = v.Target
	}

	return newTweets
}
