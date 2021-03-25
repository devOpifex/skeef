package app

import (
	"log"
	"net/http"

	"github.com/devOpifex/skeef-app/graph"
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
