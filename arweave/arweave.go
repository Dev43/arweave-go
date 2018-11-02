package arweave

import (
	"context"
	"encoding/base64"
	"math/rand"
)

func NewArweave(ctx context.Context, url string) (*ArweaveClient, error) {
	return Dial(ctx, url)
}

func (c *ArweaveClient) CreateTransaction(w *Wallet, data []byte) (*JsonTransaction, error) {
	lastTx, err := c.LastTransaction(w.Address())
	if err != nil {
		return nil, err
	}
	price, err := c.GetReward(data)
	if err != nil {
		return nil, err
	}
	id := base64.RawURLEncoding.EncodeToString(RandBytes(32))
	tx := Transaction{
		id:       id,
		lastTx:   lastTx,
		owner:    w.Public(),
		quantity: "0",
		data:     base64.RawURLEncoding.EncodeToString(data),
		reward:   price,
		txType:   "data",
	}
	tx.Sign(w)
	txa := JsonTransaction{
		Id:        id,
		LastTx:    lastTx,
		Owner:     w.Public(),
		Tags:      []interface{}{},
		Target:    "xblmNxr6cqDT0z7QIWBCo8V0UfJLd3CRDffDhF5Uh9g",
		Quantity:  "10",
		Data:      base64.RawURLEncoding.EncodeToString(data),
		Reward:    "2000",
		Signature: tx.signature,
		// TxType:    "transfer",
	}
	return &txa, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
