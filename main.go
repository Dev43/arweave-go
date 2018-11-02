package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Dev43/arweave-go/arweave"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ar, err := arweave.NewArweave(context.TODO(), "http://159.89.121.10:1984")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	// a, err := ar.GetData("qbQPYgZb3q_XHOBUDkKNCxPGMAf9tH4XIwasGBrp1mY")
	// fmt.Println(a)
	// fmt.Println(err)
	// b, err := ar.GetTransaction("LY4cgty-7IR-bB7J0aLcge4c6NZJlBMVhDJLYROJp1o")

	// fmt.Printf("%+v\n", b)
	// fmt.Println(err)
	// c, err := ar.LastTransaction("aE1AjkBoXBfF-PRP2dzRrbYY8cY2OYzeH551nSPRU5M")
	// fmt.Println(c)
	// fmt.Println(err)
	// d, err := ar.GetInfo()
	// fmt.Println(d)
	// fmt.Println(err)
	// e, err := ar.GetTransactionField("LY4cgty-7IR-bB7J0aLcge4c6NZJlBMVhDJLYROJp1o", "id")
	// fmt.Println(e)
	// fmt.Println(err)

	// e, err = ar.GetTransactionField("LY4cgty-7IR-bB7J0aLcge4c6NZJlBMVhDJLYROJp1o", "data.html")
	// fmt.Println(e)
	// fmt.Println(err)

	// f, err := ar.GetBlockByID("3obrAN4pwwkO9HSGCjW3AOav3cjUIwcP3Ewo_ev7W9GLJeeW69qVeAvdelOVlA7c")
	// fmt.Println(f)
	// fmt.Println(err)
	// g, err := ar.GetBlockByHeight(0)
	// fmt.Println(g)
	// fmt.Println(err)
	// h, err := ar.GetCurrentBlock()
	// fmt.Println(h)
	// fmt.Println(err)
	// e, err = ar.GetBalance("aE1AjkBoXBfF-PRP2dzRrbYY8cY2OYzeH551nSPRU5M")
	// fmt.Println(e)
	// fmt.Println(err)
	// j, err := ar.GetPeers()
	// fmt.Println(j)
	// fmt.Println(err)

	w := arweave.NewWallet()
	err = w.ExtractKey("arweave.json")
	if err != nil {
		log.Fatal(err)
	}
	tx, err := ar.CreateTransaction(w, []byte(""))
	// sig, err := w.Sign(nil)
	if err != nil {
		log.Fatal(err)
	}
	bb, err := json.Marshal(tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bb))

	// resp, err := ar.Commit(bb)
	// fmt.Println(resp)
	// fmt.Println(err)
}
