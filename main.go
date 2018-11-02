package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dev43/arweave-go/arweave"
)

func main() {
	ar, err := arweave.NewArweave(context.TODO(), "http://127.0.0.1:1984")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")
	w := arweave.NewWallet()
	err = w.ExtractKey("arweave.json")
	if err != nil {
		log.Fatal(err)
	}
	// Create a basic transaction (hardcoded for now)
	tx, err := ar.CreateTransaction(w, []byte(""))
	if err != nil {
		log.Fatal(err)
	}

	serialized, err := json.Marshal(tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(serialized))

	resp, err := ar.Commit(serialized)
	fmt.Println(resp)
	fmt.Println(err)
}
