package app

import (
	"encoding/json"
	"fmt"

	"github.com/devOpifex/skeef-app/graph"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/websocket"
)

func (app *Application) readSocket(con *websocket.Conn) {
	for {
		msgType, msg, err := con.ReadMessage()

		if err != nil {
			fmt.Println(err)
			return
		}

		app.InfoLog.Printf("%s\n", string(msg))

		if !app.Database.StreamOnGoing() {
			continue
		}

		tokens, err := app.Database.GetTokens()

		if err != nil {
			continue
		}

		search := app.Database.GetActiveStream()
		var index int

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
		var stream, _ = client.Streams.Filter(params)

		var demux = twitter.NewSwitchDemux()
		demux.Tweet = func(tweet *twitter.Tweet) {
			index++
			fmt.Printf("Tweet #%v\n", index)
			nodes, edges := graph.GetUserNet(*tweet)
			hashNodes, hashEdges := graph.GetHashNet(*tweet)
			nodes = append(nodes, hashNodes...)
			edges = append(edges, hashEdges...)
			g := graph.BuildGraph(nodes, edges)
			j, _ := json.Marshal(g)
			con.WriteMessage(msgType, []byte(j))
		}

		for message := range stream.Messages {
			demux.Handle(message)
		}

	}
}
