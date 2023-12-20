package main

import (
	"flag"
	"fmt"
)

func calcPrice(btcPrice, flo, paz, boris float64) {
	fmt.Printf("Flo:\t %.8f BTC\n", flo/btcPrice)
	fmt.Printf("Paz:\t %.8f BTC\n", paz/btcPrice)
	fmt.Printf("Boris:\t %.8f BTC\n", boris/btcPrice)
}

func main() {
	var flo, paz, boris, btcPrice float64

	flag.Float64Var(&flo, "flo", 15.0, "Flo's share. Default: 15.0")
	flag.Float64Var(&paz, "paz", 40.0, "Paz's share. Default: 40.0")
	flag.Float64Var(&boris, "boris", 40.0, "Boris' share. Default: 40.0")

	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Usage: calc -flo <flo_value> -paz <paz_value> -boris <boris_value> <btc_price>")
		return
	}

	btcPriceInput := args[0]
	_, err := fmt.Sscanf(btcPriceInput, "%f", &btcPrice)

	if err != nil {
		fmt.Println("Error parsing BTC price:", err)
		return
	}

	calcPrice(btcPrice, flo, paz, boris)
}
