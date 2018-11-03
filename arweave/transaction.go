package arweave

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

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
func (t *Transaction) Id() [32]byte {
	return t.id
}

type signature struct {
	Signature string `json:"signature"`
}

func (t *Transaction) Sign(w *Wallet) error {
	payload := []byte(t.formatMsg())
	msg := sha256.Sum256(payload)
	sig, err := w.Sign(msg[:])
	if err != nil {
		return err
	}
	// Ensure sig is valid
	err = w.Verify(msg[:], sig)
	if err != nil {
		return err
	}

	// Now let's calculate the id
	// Note the arweave client takes the SHA256 of the base64url encoded signature
	id := sha256.Sum256([]byte(base64.RawURLEncoding.EncodeToString(sig)))

	// add them to our transaction
	t.signature = sig
	t.id = id
	return nil
}

// Creates the message that needs to be signed
func (t *Transaction) formatMsg() string {
	return fmt.Sprintf("%s%s%s%s%s%s", t.Owner(), t.Target(), t.Data(), t.Quantity(), t.Reward(), t.LastTx())
}

func (t *Transaction) FormatJson() *JsonTransaction {
	// base64url Encode all the things
	return &JsonTransaction{
		Id:        base64.RawURLEncoding.EncodeToString(t.id[:]),
		LastTx:    base64.RawURLEncoding.EncodeToString([]byte(t.lastTx)),
		Owner:     base64.RawURLEncoding.EncodeToString([]byte(t.owner)),
		Tags:      t.tags,
		Target:    base64.RawURLEncoding.EncodeToString([]byte(t.target)),
		Quantity:  t.quantity,
		Data:      base64.RawURLEncoding.EncodeToString([]byte(t.data)),
		Reward:    t.reward,
		Signature: base64.RawURLEncoding.EncodeToString(t.signature),
	}
}
