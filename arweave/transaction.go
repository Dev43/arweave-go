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

func (t *Transaction) Owner() []byte {
	return t.owner.Bytes()
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
func (t *Transaction) Tags() []map[string]interface{} {
	return t.tags
}

type signature struct {
	Signature string `json:"signature"`
}

func (t *Transaction) Sign(w *Wallet) error {
	payload := t.formatMsgBytes()
	// payload := []byte(t.formatMsg())
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
	id := sha256.Sum256((sig))

	// add them to our transaction
	t.signature = sig
	t.id = id
	return nil
}

// Creates the message that needs to be signed
func (t *Transaction) formatMsg() string {
	return fmt.Sprintf("%s%s%s%s%s%s", t.Owner(), t.Target(), t.Data(), t.Quantity(), t.Reward(), t.LastTx(), t.Tags())
}

func (t *Transaction) formatMsgBytes() []byte {
	var msg []byte
	lastTx, err := base64.RawURLEncoding.DecodeString(t.LastTx())
	if err != nil {
		fmt.Println("err", err)
	}
	target, err := base64.RawURLEncoding.DecodeString(t.Target())
	if err != nil {
		fmt.Println("err", err)
	}
	msg = append(msg, []byte(t.Owner())...)
	msg = append(msg, target...)
	msg = append(msg, []byte(t.Data())...)
	msg = append(msg, []byte(t.Quantity())...)
	msg = append(msg, []byte(t.Reward())...)
	msg = append(msg, lastTx...)

	return msg
}

func (t *Transaction) FormatJson() *JsonTransaction {
	// base64url Encode all the things
	return &JsonTransaction{
		Id:     base64.RawURLEncoding.EncodeToString(t.id[:]),
		LastTx: (t.lastTx),
		// LastTx: base64.RawURLEncoding.EncodeToString([]byte(t.lastTx)),
		// Owner:    t.owner,
		Owner:    base64.RawURLEncoding.EncodeToString([]byte(t.owner.Bytes())),
		Tags:     t.tags,
		Target:   (t.target),
		Quantity: t.quantity,
		Data:     base64.RawURLEncoding.EncodeToString([]byte("")),
		// Reward:    t.reward,
		Reward:    "3211792120",
		Signature: base64.RawURLEncoding.EncodeToString(t.signature),
	}
}
