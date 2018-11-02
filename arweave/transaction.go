package arweave

import (
	"encoding/base64"
	"fmt"
)

//A wallet address is a base64url encoded SHA256 hash of the raw unencoded RSA modulus.
type Transaction struct {
	id        string // A SHA2-256 hash of the signature, based 64 URL encoded.
	lastTx    string // The ID of the last transaction made from the same address base64url encoded. If no previous transactions have been made from the address this field is set to an empty string.
	owner     string //The modulus of the RSA key pair corresponding to the wallet making the transaction, base64url encoded.
	target    string //  If making a financial transaction this field contains the wallet address of the recipient base64url encoded. If the transaction is not a financial this field is set to an empty string.
	quantity  string // If making a financial transaction this field contains the amount in Winston to be sent to the receiving wallet. If the transaction is not financial this field is set to the string "0". 1 AR = 1000000000000 (1e+12) Winston
	txType    string // If making a financial transaction this field contains the amount in Winston to be sent to the receiving wallet. If the transaction is not financial this field is set to the string "0". 1 AR = 1000000000000 (1e+12) Winston
	data      string //If making an archiving transaction this field contains the data to be archived base64url encoded. If the transaction is not archival this field is set to an empty string.
	reward    string //  This field contains the mining reward for the transaction in Winston.
	signature string // The data for the signature is comprised of previous data from the rest of the transaction.
}

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
		txType:   "data",
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
	return fmt.Sprintf("%s%s%s%s%s%s%s", t.Owner(), t.Target(), t.Id(), t.Data(), t.Quantity(), t.Reward(), t.LastTx())
}
