package main

import (
	"fmt"
	"log"

	"github.com/Dev43/arweave-go/arweave"
)

func main() {
	// create a new arweave client
	ar, err := arweave.NewArweaveClient("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	// create a new wallet instance
	w := arweave.NewWallet()
	// extract the key from the wallet instance
	err = w.ExtractKey("arweave.json")
	if err != nil {
		log.Fatal(err)
	}
	// create a transaction
	tx, err := ar.CreateTransaction(w, "0", "I am on the weave", "xblmNxr6cqDT0z7QIWBCo8V0UfJLd3CRDffDhF5Uh9g")
	if err != nil {
		log.Fatal(err)
	}

	// sign the transaction
	err = tx.Sign(w)
	if err != nil {
		log.Fatal(err)
	}

	// send the transaction
	resp, err := ar.SendTransaction(tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	txn := tx.Format()
	pendingTx, err := ar.GetTransaction(txn.Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pendingTx)
}
