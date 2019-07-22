# Arweave Go SDK


[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Dev43/arweave-go/blob/master/LICENSE.md)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/Dev43/arweave-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dev43/arweave-go)](https://goreportcard.com/report/github.com/Dev43/arweave-go)

Golang SDK for the Arweave client.

Example of use:

```golang

	// create a new arweave client
    ar, err := transactor.NewTransactor("209.97.142.169")
    if err != nil {
        log.Fatal(err)
    }

	// create a new wallet instance
	w := wallet.NewWallet()
	// extract the key from the wallet instance
	err = w.ExtractKey("./arweave.json")
	if err != nil {
		log.Fatal(err)
	}
	// create a transaction
	txBuilder, err := ar.CreateTransaction(context.TODO(), w, "0", []byte(""), "xblmNxr6cqDT0z7QIWBCo8V0UfJLd3CRDffDhF5Uh9g")
	if err != nil {
		log.Fatal(err)
	}

	// sign the transaction
	txn, err := txBuilder.Sign(w)
	if err != nil {
		log.Fatal(err)
	}

	// send the transaction
	resp, err := ar.SendTransaction(context.TODO(), txn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)

	// wait for the transaction to get mined
	pendingTx, err := ar.WaitMined(context.TODO(), txn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pendingTx)

```
