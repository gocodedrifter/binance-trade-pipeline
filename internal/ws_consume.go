package internal

import (
	"fmt"

	"github.com/gocodedrifter/coinpanel/pkg/websocket"
)

// WsServe : serving value from websocket
func WsServe(stopStreaming chan bool, symbolC chan string) {
	for {
		select {
		case symbol := <-symbolC:
			wsTradeHandler := func(event *websocket.WsTradeEvent) {
				fmt.Println(event)
			}
			errHandler := func(err error) {
				fmt.Println(err)
			}
			doneC, stopC, err := websocket.WsTradeServe(symbol, wsTradeHandler, errHandler)
			if err != nil {
				fmt.Println(err)
				return
			}

			go func() {
				select {
				case ok := <-stopStreaming:
					if ok {
						close(stopC)
					}
				case <-doneC:
				}
			}()
		}
	}
}
