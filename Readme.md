**Ingestion Pipeline**

**Tech Stack :** 
1. Go 1.12
2. mongodb (replica set)

**Details :** 
Data from binance websocket is saved in mongodb. 
the program listen to change stream from mongodb if new trade data is saved.

**How to run the program :**
1. make sure docker installed
2. run : **docker-compose run --rm app**

