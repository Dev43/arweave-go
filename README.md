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

### Api

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

### Batching, chunking and recombining

Arweave currently has a 3MB transaction limit. To let you be able to upload more than 3MB of data, this library lets you chunk your data into multiples chunks, and upload these chunks to the weave. These chunks are all backlinked to their previous chunk (except the first one) so all you need to use to retrieve all the chunks is the last chunk's address (the "tip" of this linked list).

So far the library has been tested with Tarred and Gzipped files (it's a good idea to compress your files, you'll be paying less fees to the network!)

#### Creating chunks

To create chunks, use the `chunker` package. Currently the chunker takes 500KB chunks, encodes it to the same encoding as the Arweave (Raw Base64 URL encoding) and adds it to a transaction.

```golang
	f, err := os.Open("example.tar.gz")
	if err != nil {
		//...
	}
	info, err := f.Stat()
	if err != nil {
		//...
	}
	chunker, err := chunker.NewChunker(f, info.Size())
	if err != nil {
		//...
	}
	chunks, err := chunker.ChunkAll()
	if err != nil {
		//...
	}
```

#### Sending a batch of transactions

To directly do both the chunking and the batch sending of transaction, you can use the `batchchunker` package. The package needs to wait for a single transaction to be mined  before sending a new one. This will take a significant amount of time as the Arweave has a ~2 minute block time.

```golang
	f, err := os.Open("example.tar.gz")
	if err != nil {
		//...
	}
	info, err := f.Stat()
	if err != nil {
		//...
	}

	// creates a new batch
	newB := batchchunker.NewBatch(ar, w, f, info.Size())

	// sends all the transactions
	list, err := newB.SendBatchTransaction()
	if err != nil {
		//...
	}
```

#### Recombining

Recombining all the chunks needs to be done from the tip of the batch chain (the last transaction sent). The function will grab all the chunks and recombine it into an io.Writer.

```golang
	// in this example, we will save the data as a file
	f, err := os.Create("example.tar.gz")
	if err != nil {
		//...
	}
	// create a new combiner
	newBCombiner := combiner.NewBatchCombiner(ar.Client)

	// grab all the chunks starting from chunk at this address
	liveChunks, err := newBCombiner.GetAllChunks("fBwTnPGNSAtHTNr6pGvqiYpLmBvTjjJOu3kBxuuit1c")
	if err != nil {
		//...
	}
	// recombine all the chunks into our file
	err = combiner.Recombine(liveChunks, f)
	if err != nil {
		//...
	}
```

If you actually run the last example, it will retrieve a folder with the Bitcoin, Ethereum and Arweave whitepaper


If you enjoy the library, please consider donating:

- *Arweave Address*: `pfJXiTwwjQwSJF9VT1ZK6kauvobWuKKLUzjz29R1gbQ`
- *Ethereum Address*: `0x3E42b8b399dca71b5c004921Fc6eFfa8dDc9409d`