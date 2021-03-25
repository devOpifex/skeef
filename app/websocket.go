package app

import (
	"encoding/json"
	"fmt"
	"io"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (app *Application) wsUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}

func (app *Application) wsReader(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		app.InfoLog.Printf("%s\n", string(msg))

		if !app.Database.StreamOnGoing() {
			if app.Stream != nil {
				app.Stream.Stop()
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
			j, _ := json.Marshal(m)
			conn.WriteMessage(messageType, []byte(j))
		}

		for message := range app.Stream.Messages {
			demux.Handle(message)
		}
	}
}

func (app *Application) wsWriter(conn *websocket.Conn) {
	for {
		fmt.Println("Sending")
		messageType, r, err := conn.NextReader()
		if err != nil {
			fmt.Println(err)
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			fmt.Println(err)
			return
		}
		if err := w.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}
}
