package arweave

import (
	"encoding/base64"
	"fmt"
)

// Really change all the data in a txn to be private, and make getters for it.

func NewTransaction(lastTx string, owner string, quantity string, txType string, data string, reward string, id string, target string, tags []string) *Transaction {
	return &Transaction{
		id:       id,
		lastTx:   lastTx,
		owner:    owner,
		quantity: quantity,
		data:     data,
		reward:   reward,
		target:   target,
	}
}

func (t *Transaction) Data() string {
	return t.data
}
func (t *Transaction) LastTx() string {
	return t.lastTx
}

func (t *Transaction) Owner() string {
	return t.owner
}
func (t *Transaction) Quantity() string {
	return t.quantity
}

func (t *Transaction) Reward() string {
	return t.reward
}

func (t *Transaction) Target() string {
	return t.target
}
func (t *Transaction) Id() string {
	return t.id
}

func (t *Transaction) Sign(w *Wallet) error {
	payload := []byte(t.formatMsg())
	sig, err := w.Sign(payload)
	if err != nil {
		return err
	}
	// Ensure sig is valid
	err = w.Verify(payload, sig)
	if err != nil {
		return err
	}
	// t.signature = sig
	t.signature = base64.RawURLEncoding.EncodeToString([]byte(sig))
	return nil
}

// Creates the message that needs to be signed
func (t *Transaction) formatMsg() string {
	return fmt.Sprintf("%s%s%s%s%s%s", t.Owner(), t.Target(), t.Data(), t.Quantity(), t.Reward(), t.LastTx())
}

func (t *Transaction) FormatJson() *JsonTransaction {
	return &JsonTransaction{
		Id:        t.id,
		LastTx:    t.lastTx,
		Owner:     t.owner,
		Tags:      t.tags,
		Target:    t.target,
		Quantity:  t.quantity,
		Data:      t.data,
		Reward:    t.reward,
		Signature: t.signature,
	}
}
