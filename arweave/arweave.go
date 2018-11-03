package arweave

import (
	"context"
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
	// Non encoded transaction fields
	tx := Transaction{
		lastTx:   lastTx,
		owner:    w.publicKey,
		quantity: "100000000000",
		target:   "xblmNxr6cqDT0z7QIWBCo8V0UfJLd3CRDffDhF5Uh9g",
		data:     "",
		reward:   "321179212",
		tags:     make([]map[string]interface{}, 0),
	}

	err = tx.Sign(w)
	if err != nil {
		return nil, err
	}

	return tx.FormatJson(), nil
}
