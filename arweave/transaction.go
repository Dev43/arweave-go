package arweave

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Data returns the data of the transaction
func (t *Transaction) Data() string {
	return t.data
}

// LastTx returns the last transaction of the account
func (t *Transaction) LastTx() string {
	return t.lastTx
}

// Owner returns the Owner of the transaction
func (t *Transaction) Owner() []byte {
	return t.owner.Bytes()
}

// Quantity returns the quantity of the transaction
func (t *Transaction) Quantity() string {
	return t.quantity
}

// Reward returns the reward of the transaction
func (t *Transaction) Reward() string {
	return t.reward
}

// Target returns the target of the transaction
func (t *Transaction) Target() string {
	return t.target
}

// Id returns the id of the transaction which is the SHA256 of the signature
func (t *Transaction) Id() [32]byte {
	return t.id
}

// Tags returns the tags of the transaction
func (t *Transaction) Tags() []map[string]interface{} {
	return t.tags
}

// Sign creates the signing message, and signs it using the private key,
// It takes the SHA256 of the resulting signature to calculate the id of
// the signature
func (t *Transaction) Sign(w *Wallet) error {
	// format the message
	payload, err := t.formatMsgBytes()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// take the SHA256 of it
	msg := sha256.Sum256(payload)

	// sign it using the RSA private key
	sig, err := w.Sign(msg[:])
	if err != nil {
		return err
	}

	// ensure the signature is valid
	err = w.Verify(msg[:], sig)
	if err != nil {
		return err
	}

	// calculate the transaction id
	id := sha256.Sum256((sig))

	// add them to our transaction
	t.signature = sig
	t.id = id
	return nil
}

// formatMsgBytes formats the message that needs to be signed. All fields
// need to be an array of bytes originating from the necessary data (not base64url encoded).
// The signing message is the SHA256 of the concatenation of the byte arrays
// of the owner public key, target address, data, quantity, reward and last transaction
func (t *Transaction) formatMsgBytes() ([]byte, error) {
	var msg []byte
	lastTx, err := base64.RawURLEncoding.DecodeString(t.LastTx())
	if err != nil {
		return nil, err
	}
	target, err := base64.RawURLEncoding.DecodeString(t.Target())
	if err != nil {
		return nil, err
	}

	msg = append(msg, []byte(t.Owner())...)
	msg = append(msg, target...)
	msg = append(msg, []byte(t.Data())...)
	msg = append(msg, []byte(t.Quantity())...)
	msg = append(msg, []byte(t.Reward())...)
	msg = append(msg, lastTx...)

	return msg, nil
}

// Format formats the transactions to a JSONTransaction that can be sent out to an arweave node
func (t *Transaction) Format() *JSONTransaction {
	return &JSONTransaction{
		Id:        base64.RawURLEncoding.EncodeToString(t.id[:]),
		LastTx:    (t.lastTx),
		Owner:     base64.RawURLEncoding.EncodeToString([]byte(t.owner.Bytes())),
		Tags:      t.tags,
		Target:    (t.target),
		Quantity:  t.quantity,
		Data:      base64.RawURLEncoding.EncodeToString([]byte(t.Data())),
		Reward:    t.reward,
		Signature: base64.RawURLEncoding.EncodeToString(t.signature),
	}
}
