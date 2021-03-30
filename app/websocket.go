package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devOpifex/skeef-app/graph"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/websocket"
)

type message struct {
	Graph       graph.Graph `json:"graph"`
	TweetsCount int         `json:"tweetsCount"`
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

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (c *Client) Read(app *Application) {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

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
				client.Conn.WriteJSON(message{Graph: app.Graph, TweetsCount: app.Count})
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
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (app *Application) wsUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.ErrorLog.Println(err)
		return ws, err
	}
	return ws, nil
}

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

func (app *Application) StopStream() {

	if app.Stream != nil {
		app.Stream.Stop()
	}

	app.Quit <- struct{}{}
	app.Database.PauseAllStreams()
}

func (app *Application) StartStream() {

	app.Count = 0
	app.Graph = graph.Graph{}

	tokens, err := app.Database.GetTokens()

	if err != nil {
		return
	}

	search := app.Database.GetActiveStream()

	var twitterConfig = oauth1.NewConfig(tokens.ApiKey, tokens.ApiSecret)
	var token = oauth1.NewToken(tokens.AccessToken, tokens.AccessSecret)

	var httpClient = twitterConfig.Client(oauth1.NoContext, token)

	var client = twitter.NewClient(httpClient)

	var params = &twitter.StreamFilterParams{
		Track:         []string{search.Track},
		StallWarnings: twitter.Bool(true),
	}
	app.Stream, _ = client.Streams.Filter(params)

	var demux = twitter.NewSwitchDemux()
	demux.Tweet = app.demux()

	for message := range app.Stream.Messages {
		demux.Handle(message)
	}
}

// demux Demux tweets and broadcast
func (app *Application) demux() func(tweet *twitter.Tweet) {
	return func(tweet *twitter.Tweet) {
		app.Count++
		app.InfoLog.Printf("Tweet #%v\n", app.Count)
		nodes, edges := graph.GetUserNet(*tweet)
		hashNodes, hashEdges := graph.GetHashNet(*tweet)
		nodes = append(nodes, hashNodes...)
		edges = append(edges, hashEdges...)
		for key := range edges {
			app.Graph.UpsertEdge(&edges[key])
		}
		for key := range nodes {
			app.Graph.UpsertNode(&nodes[key])
		}
		g := graph.Graph{Edges: edges, Nodes: nodes}

		m := message{
			Graph:       g,
			TweetsCount: app.Count,
		}
		app.Pool.Broadcast <- m
	}
}
