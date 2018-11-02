package arweave

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
)

func NewArweave(ctx context.Context, url string) (*ArweaveClient, error) {
	return Dial(ctx, url)
}

func (c *ArweaveClient) CreateTransaction(w *Wallet, data []byte) (*JsonTransaction, error) {
	lastTx, err := c.LastTransaction(w.Address())
	if err != nil {
		return nil, err
	}
	// price, err := c.GetReward(data)
	// if err != nil {
	// 	return nil, err
	// }
	tx := Transaction{
		lastTx:   lastTx,
		owner:    w.Public(),
		quantity: "10",
		target:   "xblmNxr6cqDT0z7QIWBCo8V0UfJLd3CRDffDhF5Uh9g",
		data:     base64.RawURLEncoding.EncodeToString(data),
		reward:   "2000",
		tags:     []interface{}{},
	}
	tx.Sign(w)
	h := sha256.New()
	h.Write([]byte(tx.signature))
	id := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	tx.id = id

	return tx.FormatJson(), nil
}
