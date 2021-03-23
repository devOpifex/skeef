package stream

import (
	"fmt"
	"time"
)

type Stream struct {
	Name      string
	Follow    string
	Track     string
	Locations string
	Active    int
}

func StartStream(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second * 3)
			fmt.Println("Running")
		}
	}
}
