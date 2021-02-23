package main

import (
	"fmt"

	"github.com/gocodedrifter/coinpanel/pkg/websocket"
)

func main() {
	wsTradeHandler := func(event *websocket.WsTradeEvent) {
		fmt.Println(event)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := websocket.WsTradeServe("LTCBTC", wsTradeHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
