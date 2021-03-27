package app

import (
	"fmt"
	"log"
	"net/http"

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

func (c *Client) Read(app *Application) {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if !app.Database.StreamOnGoing() {

			if app.Stream != nil {
				app.Stream.Stop()
				app.Stream = nil
			}

			continue
		}

		tokens, err := app.Database.GetTokens()

		if err != nil {
			continue
		}

		search := app.Database.GetActiveStream()

		var twitterConfig = oauth1.NewConfig(tokens.ApiKey, tokens.ApiSecret)
		var token = oauth1.NewToken(tokens.AccessToken, tokens.AccessSecret)

		// http.Client will automatically authorize Requests
		var httpClient = twitterConfig.Client(oauth1.NoContext, token)

		// Twitter client
		var client = twitter.NewClient(httpClient)

		var params = &twitter.StreamFilterParams{
			Track:         []string{search.Track},
			StallWarnings: twitter.Bool(true),
		}
		app.Stream, _ = client.Streams.Filter(params)

		var demux = twitter.NewSwitchDemux()
		demux.Tweet = func(tweet *twitter.Tweet) {
			app.Count++
			app.InfoLog.Printf("Tweet #%v\n", app.Count)
			nodes, edges := graph.GetUserNet(*tweet)
			hashNodes, hashEdges := graph.GetHashNet(*tweet)
			nodes = append(nodes, hashNodes...)
			edges = append(edges, hashEdges...)
			g := graph.BuildGraph(nodes, edges)
			m := message{
				Graph:       g,
				TweetsCount: app.Count,
			}
			c.Pool.Broadcast <- m
		}

		for message := range app.Stream.Messages {
			demux.Handle(message)
		}

	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			for client := range pool.Clients {
				client.Conn.WriteJSON(message{})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			for client := range pool.Clients {
				client.Conn.WriteJSON(message{})
			}
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
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

func (app *Application) socket(pool *Pool, w http.ResponseWriter, r *http.Request) {
	ws, err := app.wsUpgrade(w, r)

	if err != nil {
		log.Println(err)
	}

	client := &Client{
		Conn: ws,
		Pool: pool,
	}

	pool.Register <- client
	client.Read(app)

}
