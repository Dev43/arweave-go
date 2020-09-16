package tx

import (
	"crypto/sha256"
	"encoding/json"
	"math/big"

	"github.com/Dev43/arweave-go"
	"github.com/Dev43/arweave-go/utils"
)

// NewTransaction creates a brand new transaction struct
func NewTransaction(lastTx string, owner *big.Int, quantity string, target string, data []byte, reward string) *Transaction {
	return &Transaction{
		lastTx:   lastTx,
		owner:    owner,
		quantity: quantity,
		target:   target,
		data:     data,
		reward:   reward,
		tags:     make([]Tag, 0),
	}
}

// Data returns the data of the transaction
func (t *Transaction) Data() string {
	return utils.EncodeToBase64(t.data)
}

// RawData returns the unencoded data
func (t *Transaction) RawData() []byte {
	return t.data
}

// LastTx returns the last transaction of the account
func (t *Transaction) LastTx() string {
	return t.lastTx
}

// Owner returns the Owner of the transaction
func (t *Transaction) Owner() string {
	return utils.EncodeToBase64(t.owner.Bytes())
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
func (t *Transaction) ID() []byte {
	return t.id
}

// Hash returns the base64 RawURLEncoding of the transaction hash
func (t *Transaction) Hash() string {
	return utils.EncodeToBase64(t.id)
}

// Tags returns the tags of the transaction in plain text
func (t *Transaction) Tags() ([]Tag, error) {
	tags := []Tag{}
	for _, tag := range t.tags {
		// access name
		tagName, err := utils.DecodeString(tag.Name)
		if err != nil {
			return nil, err
		}
		tagValue, err := utils.DecodeString(tag.Value)
		if err != nil {
			return nil, err
		}
		tags = append(tags, Tag{Name: string(tagName), Value: string(tagValue)})
	}
	return tags, nil
}

// RawTags returns the unencoded tags of the transaction
func (t *Transaction) RawTags() []Tag {
	return t.tags
}

// AddTag adds a new tag to the transaction
func (t *Transaction) AddTag(name string, value string) error {
	tag := Tag{
		Name:  utils.EncodeToBase64([]byte(name)),
		Value: utils.EncodeToBase64([]byte(value)),
	}
	t.tags = append(t.tags, tag)
	return nil
}

func (t *Transaction) SetID(id []byte) {
	t.id = id
}

func (t *Transaction) SetSignature(signature []byte) {
	t.signature = signature
}

// Signature returns the signature of the transaction
func (t *Transaction) Signature() string {
	return utils.EncodeToBase64(t.signature)
}

// Sign creates the signing message, and signs it using the private key,
// It takes the SHA256 of the resulting signature to calculate the id of
// the signature
func (t *Transaction) Sign(w arweave.WalletSigner) (*Transaction, error) {
	// format the message
	payload, err := t.FormatMsgBytes()
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

	idB := make([]byte, len(id))
	copy(idB, id[:])
	t.SetID(idB)

	// we copy t into tx
	tx := Transaction(*t)
	// add the signature and ID to our new signature struct
	tx.signature = sig

	return &tx, nil
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
	id, err := utils.DecodeString(txn.ID)
	if err != nil {
		return err
	}
	t.id = id

	t.lastTx = txn.LastTx

	// gives me byte representation of the big num
	owner, err := utils.DecodeString(txn.Owner)
	if err != nil {
		return err
	}
	n := new(big.Int)
	t.owner = n.SetBytes(owner)

	t.tags = txn.Tags
	t.target = txn.Target
	t.quantity = txn.Quantity

	data, err := utils.DecodeString(txn.Data)
	if err != nil {
		return err
	}
	t.data = data
	t.reward = txn.Reward

	sig, err := utils.DecodeString(txn.Signature)
	if err != nil {
		return err
	}
	t.signature = sig

	return nil
}

// FormatMsgBytes formats the message that needs to be signed. All fields
// need to be an array of bytes originating from the necessary data (not base64url encoded).
// The signing message is the SHA256 of the concatenation of the byte arrays
// of the owner public key, target address, data, quantity, reward and last transaction
func (t *Transaction) FormatMsgBytes() ([]byte, error) {
	var msg []byte
	lastTx, err := utils.DecodeString(t.LastTx())
	if err != nil {
		return nil, err
	}
	target, err := utils.DecodeString(t.Target())
	if err != nil {
		return nil, err
	}

	tags, err := t.encodeTagData()
	if err != nil {
		return nil, err
	}

	msg = append(msg, t.owner.Bytes()...)
	msg = append(msg, target...)
	msg = append(msg, t.data...)
	msg = append(msg, t.quantity...)
	msg = append(msg, t.reward...)
	msg = append(msg, lastTx...)
	msg = append(msg, tags...)

	return msg, nil
}

// We need to encode the tag data properly for the signature. This means having the unencoded
// value of the Name field concatenated with the unencoded value of the Value field
func (t *Transaction) encodeTagData() (string, error) {
	tagString := ""
	unencodedTags, err := t.Tags()
	if err != nil {
		return "", err
	}
	for _, tag := range unencodedTags {
		tagString += tag.Name + tag.Value
	}

	return tagString, nil
}

// Format formats the transactions to a JSONTransaction that can be sent out to an arweave node
func (t *Transaction) format() *transactionJSON {
	return &transactionJSON{
		ID:        utils.EncodeToBase64(t.id),
		LastTx:    t.lastTx,
		Owner:     utils.EncodeToBase64(t.owner.Bytes()),
		Tags:      t.tags,
		Target:    t.target,
		Quantity:  t.quantity,
		Data:      utils.EncodeToBase64(t.data),
		Reward:    t.reward,
		Signature: utils.EncodeToBase64(t.signature),
	}
}
