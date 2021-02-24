package main

import (
	"fmt"

	"github.com/gocodedrifter/coinpanel/internal/service"
)

func main() {
	s, p := display()
	service.WSServe(s)
	service.StreamMongo(s, p)
	fmt.Scanln()
}

func display() (symbol string, price string) {
	fmt.Printf("%-25s", "Please Enter symbol : ")
	fmt.Scan(&symbol)

	fmt.Printf("%-40s", "Please Enter limit price to monitor : ")
	fmt.Scan(&price)
	return
}
