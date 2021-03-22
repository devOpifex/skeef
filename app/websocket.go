package app

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func (app *Application) readSocket(con *websocket.Conn) {
	for {
		msgType, msg, err := con.ReadMessage()

		if err != nil {
			fmt.Println(err)
			return
		}

		app.InfoLog.Printf("%s\n", string(msgType))
		app.InfoLog.Printf("%s\n", string(msg))

	}
}
