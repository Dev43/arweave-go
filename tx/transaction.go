package tx

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"github.com/Dev43/arweave-go"

)

func NewTransaction(lastTx string, owner *big.Int, quantity string, target string, data []byte, reward string, tags []map[string]interface{}) *Transaction {
	return &Transaction{
		lastTx:   lastTx,
		owner:    owner,
		quantity: quantity,
		target:   target,
		data:     data,
		reward:   reward,
		tags:     make([]map[string]interface{}, 0),
	}
}

// Data returns the data of the transaction
func (t *Transaction) Data() string {
	return base64.RawURLEncoding.EncodeToString(t.data)
}

// LastTx returns the last transaction of the account
func (t *Transaction) LastTx() string {
	return t.lastTx
}

// Owner returns the Owner of the transaction
func (t *Transaction) Owner() string {
	return base64.RawURLEncoding.EncodeToString(t.owner.Bytes())
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

// ID returns the id of the transaction which is the SHA256 of the signature
func (t *Transaction) ID() [32]byte {
	return t.id
}

// Hash returns the base64 RawURLEncoding of the transaction hash
func (t *Transaction) Hash() string {
	return base64.RawURLEncoding.EncodeToString(t.id[:])
}

// Tags returns the tags of the transaction
func (t *Transaction) Tags() []map[string]interface{} {
	return t.tags
}

// Signature returns the signature of the transaction
func (t *Transaction) Signature() string {
	return base64.RawURLEncoding.EncodeToString(t.signature)
}

// Sign creates the signing message, and signs it using the private key,
// It takes the SHA256 of the resulting signature to calculate the id of
// the signature
func (t *Transaction) Sign(w arweave.WalletSigner) (*Transaction, error) {
	// format the message
	payload, err := t.formatMsgBytes()
	if err != nil {
		return nil, err
	}

	msg := sha256.Sum256(payload)

	sig, err := w.Sign(msg[:])
	if err != nil {
		return nil, err
	}

	err = w.Verify(msg[:], sig)
	if err != nil {
		return nil, err
	}

	id := sha256.Sum256((sig))

	// we copy t into tx
	tx := Transaction(*t)
	// add the signature and ID to our new signature struct
	tx.signature = sig
	tx.id = id

	return &tx, nil
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

	msg = append(msg, t.owner.Bytes()...)
	msg = append(msg, target...)
	msg = append(msg, t.data...)
	msg = append(msg, t.quantity...)
	msg = append(msg, t.reward...)
	msg = append(msg, lastTx...)

	return msg, nil
}

// Format formats the transactions to a JSONTransaction that can be sent out to an arweave node
func (t *Transaction) format() *transactionJSON {
	return &transactionJSON{
		ID:        base64.RawURLEncoding.EncodeToString(t.id[:]),
		LastTx:    t.lastTx,
		Owner:     base64.RawURLEncoding.EncodeToString(t.owner.Bytes()),
		Tags:      t.tags,
		Target:    t.target,
		Quantity:  t.quantity,
		Data:      base64.RawURLEncoding.EncodeToString(t.data),
		Reward:    t.reward,
		Signature: base64.RawURLEncoding.EncodeToString(t.signature),
	}
}

// MarshalJSON marshals as JSON
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.format())
}

// UnmarshalJSON unmarshals as JSON
func (t *Transaction) UnmarshalJSON(input []byte) error {
	txn := transactionJSON{}
	err := json.Unmarshal(input, &txn)
	if err != nil {
		return err
	}
	id, err := base64.RawURLEncoding.DecodeString(txn.ID)
	if err != nil {
		return err
	}
	var id32 [32]byte
	copy(id32[:], id)
	t.id = id32

	t.lastTx = txn.LastTx

	// gives me byte representation of the big num
	owner, err := base64.RawURLEncoding.DecodeString(txn.Owner)
	if err != nil {
		return err
	}
	n := new(big.Int)
	t.owner = n.SetBytes(owner)

	t.tags = txn.Tags
	t.target = txn.Target
	t.quantity = txn.Quantity

	data, err := base64.RawURLEncoding.DecodeString(txn.Data)
	if err != nil {
		return err
	}
	t.data = data
	t.reward = txn.Reward

	sig, err := base64.RawURLEncoding.DecodeString(txn.Signature)
	if err != nil {
		return err
	}
	t.signature = sig

	return nil
}
