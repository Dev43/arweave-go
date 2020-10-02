# Arweave Go SDK


[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Dev43/arweave-go/blob/master/LICENSE.md)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/Dev43/arweave-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dev43/arweave-go)](https://goreportcard.com/report/github.com/Dev43/arweave-go)

Golang Client to interact with the Arweave Blockchain.

## Usage

### Wallet

In the current version, you can load the Arweave wallet file created from the Arweave server or the plugin.

```golang
// create a new wallet instance
w := wallet.NewWallet()
// extract the key from the wallet instance
err = w.LoadKeyFromFile("./arweave.json")
if err != nil {
	//...
}
```

You can directly load the key by using it's filepath or pass it as an array of bytes using `LoadKey([]byte)`.

With the wallet struct, you can sign and verify a message:

```golang
// sign the message "example"
msg := []byte("example")
sig, err := w.Sign(msg))
if err != nil {
	//...
}
err = w.Verify(msg, sig)
if err != nil {
	// message signature is not valid...
}
// message signature is valid
```

### API

You can call all of the Arweave HTTP api endpoints using the api package. First you must give it the node url or IP address.

```golang
ipAddress := "127.0.0.1"
c, err := api.Dial(ipAddress)
if err != nil {
	// problem connecting
}
```

To call the endpoints, you will need to pass in a context.

```golang
c.GetBalance(context.TODO(), "1seRanklLU_1VTGkEk7P0xAwMJfA7owA1JHW5KyZKlY")
```

### Transactions

To create a new transaction, you will need to interact with the `transactor` package. The transactor package has 3 main functions, creating, sending and waiting for a transaction.

```golang
	// create a new transactor client
	ar, err := transactor.NewTransactor("127.0.0.1")
	if err != nil {
		//...
	}

	// create a new wallet instance
	w := wallet.NewWallet()
	// extract the key from the wallet instance
	err = w.LoadKeyFromFile("./arweave.json")
	if err != nil {
		//...
	}
	// create a transaction
	txBuilder, err := ar.CreateTransaction(context.TODO(), w, "0", []byte(""), "1seRanklLU_1VTGkEk7P0xAwMJfA7owA1JHW5KyZKlY")
	if err != nil {
		//...
	}
	
	// sign the transaction
	txn, err := txBuilder.Sign(w)
	if err != nil {
		//...
	}

	// send the transaction
	resp, err := ar.SendTransaction(context.TODO(), txn)
	if err != nil {
		//...
	}

	// wait for the transaction to get mined
	finalTx, err := ar.WaitMined(context.TODO(), txn)
	if err != nil {
		//...
	}
	// get the hash of the transaction
	fmt.Println(finalTx.Hash())
```


If you enjoy the library, please consider donating:

- **Arweave Address**: `pfJXiTwwjQwSJF9VT1ZK6kauvobWuKKLUzjz29R1gbQ`
- **Ethereum Address**: `0x3E42b8b399dca71b5c004921Fc6eFfa8dDc9409d`