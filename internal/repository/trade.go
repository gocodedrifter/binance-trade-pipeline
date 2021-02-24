package trade

import (
	"context"
	"log"
	"time"

	"github.com/gocodedrifter/coinpanel/pkg/db"
	"github.com/gocodedrifter/coinpanel/pkg/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

const trade = "trade"

func getTableName() string {
	return trade
}

// Save : save json
func Save(json []byte) {
	var dataToSave interface{}
	bson.UnmarshalExtJSON(json, true, &dataToSave)
	collection := db.GetDB().Collection(getTableName())
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, &dataToSave)
	if err != nil {
		log.Println(err.Error())
	}
}

// FindData : find data based on symbol
func FindData(symbol string) (event *websocket.WsTradeEvent, err error) {
	collection := db.GetDB().Collection(getTableName())
	filter := bson.M{"s": symbol}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err = collection.FindOne(ctx, filter).Decode(&event); err != nil {
		log.Println("error : ", err.Error())
	}
	return
}

// FindAll : find all data based on symbol
func FindAll(symbol string, price float64) (events []*websocket.WsTradeEvent, err error) {
	collection := db.GetDB().Collection(getTableName())
	filter := bson.M{"s": symbol, "pnum": bson.M{"$gt": price}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	curr, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("error : ", err.Error())
	}

	for curr.Next(ctx) {
		var event *websocket.WsTradeEvent
		if err = curr.Decode(&event); err != nil {
			log.Println("error due to : ", err.Error())
		}

		events = append(events, event)
	}
	return
}
