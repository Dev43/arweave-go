# Arweave Go SDK


[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Dev43/arweave-go/blob/master/LICENSE.md)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/Dev43/arweave-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dev43/arweave-go)](https://goreportcard.com/report/github.com/Dev43/arweave-go)

Golang SDK for the Arweave client.

Example of use:

```golang
package main

import (
	"fmt"
	"log"

	"github.com/Dev43/arweave-go/transactor"
	"github.com/Dev43/arweave-go/wallet"
)

func main() {
	// create a new arweave client
	ar, err := transactor.NewTransactor("209.97.142.169")
	if err != nil {
		log.Fatal(err)
	}

	// create a new wallet instance
	w := wallet.NewWallet()
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
	pendingTx, err := ar.Client.GetTransaction(txn.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pendingTx)
}

```