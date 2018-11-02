package arweave

import (
	"encoding/base64"
	"encoding/json"
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
	data, _ := base64.RawURLEncoding.DecodeString(t.data)
	return string(data)
}
func (t *Transaction) LastTx() string {
	lastTx, _ := base64.RawURLEncoding.DecodeString(t.lastTx)
	return string(lastTx)
}

func (t *Transaction) Owner() string {
	data, _ := base64.RawURLEncoding.DecodeString(t.owner)
	return string(data)
}
func (t *Transaction) Quantity() string {
	return t.quantity
}

func (t *Transaction) Reward() string {
	return t.reward
}

func (t *Transaction) Target() string {
	data, _ := base64.RawURLEncoding.DecodeString(t.target)
	return string(data)
}
func (t *Transaction) Id() string {
	return t.id
}

type signature struct {
	Signature string `json:"signature"`
}

func (t *Transaction) Sign(w *Wallet) error {
	payload := []byte(t.formatMsg())
	fullSig, err := w.Sign(payload)
	if err != nil {
		return err
	}
	// Ensure sig is valid
	err = w.Verify(payload, fullSig)
	if err != nil {
		return err
	}
	// Extract the signature
	sig := signature{}
	json.Unmarshal([]byte(fullSig), &sig)
	// Encode it to base64url
	t.signature = base64.RawURLEncoding.EncodeToString([]byte(sig.Signature))
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
