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

		app.InfoLog.Printf("%v\n", msgType)
		app.InfoLog.Printf("%s\n", string(msg))

	}
}
