package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	trade "github.com/gocodedrifter/coinpanel/internal/repository"
	"github.com/gocodedrifter/coinpanel/pkg/db"
	"github.com/gocodedrifter/coinpanel/pkg/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const tradeDoc = "trade"

func getTableName() string {
	return tradeDoc
}

// WSServe : serve web socket based on symbol
func WSServe(symbol string) {
	go func() {
		wsTradeHandler := func(event *websocket.WsTradeEvent) {
			n, err := strconv.ParseFloat(event.Price, 64)
			if err != nil {
				log.Println("error due to : ", err.Error())
			}
			event.PriceNum = n
			mEvent, _ := json.Marshal(event)
			trade.Save(mEvent)
		}
		errHandler := func(err error) {
			log.Println(err)
		}
		doneC, _, err := websocket.WsTradeServe(symbol, wsTradeHandler, errHandler)
		if err != nil {
			log.Println(err)
			return
		}
		<-doneC
	}()
}

func iterateChangeStream(routineCtx context.Context, waitGroup sync.WaitGroup, stream *mongo.ChangeStream) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()
	for stream.Next(routineCtx) {
		changeDoc := struct {
			FullDocument map[string]interface{} `bson:"fullDocument"`
		}{}
		if err := stream.Decode(&changeDoc); err != nil {
			panic(err)
		}
		event := changeDoc.FullDocument["e"]
		symbol := changeDoc.FullDocument["s"]
		price := changeDoc.FullDocument["p"]
		quantity := changeDoc.FullDocument["q"]
		fmt.Printf("%-15s%-15s%-15s%-15s\n", event, symbol, price, quantity)
	}
}

// StreamMongo : stream data from mongo
func StreamMongo(symbol string, price string) {
	p, _ := strconv.ParseFloat(price, 64)
	fmt.Printf("%-15s%-15s%-15s%-15s\n", "Event", "Symbol", "Price", "Quantity")
	episodesCollection := db.GetDB().Collection(getTableName())
	var waitGroup sync.WaitGroup
	matchPipeline := bson.D{
		{
			"$match", bson.D{
				{"fullDocument.s", symbol},
				{"fullDocument.pnum", bson.D{
					{"$gt", p},
				}},
			},
		},
	}
	episodesStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{matchPipeline})
	if err != nil {
		panic(err)
	}
	waitGroup.Add(1)
	routineCtx, _ := context.WithCancel(context.Background())
	go iterateChangeStream(routineCtx, waitGroup, episodesStream)
	waitGroup.Wait()
}
