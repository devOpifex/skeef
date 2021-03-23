package app

import (
	"fmt"
	"time"
)

func (app *Application) StartStream() {
	for {
		select {
		case <-app.StopStream:
			return
		default:
			time.Sleep(time.Second * 3)
			fmt.Println("Running")
		}
	}
}
