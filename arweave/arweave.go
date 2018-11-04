package arweave

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// DefaultURL is the local host url
const DefaultURL = "http://127.0.0.1:1984"

// DefaultPort of the arweave client
const DefaultPort = "1984"

// NewArweaveClient creates a new arweave client. You need to pass in a context and a url
// If sending an empty string, the default url is localhosts
func NewArweaveClient(fullUrl string) (*ArweaveClient, error) {
	furl := fullUrl
	if furl == "" {
		return Dial(DefaultURL)
	}
	u, err := url.Parse(furl)
	if err != nil {
		return nil, err
	}
	if u.Port() == "" {
		furl = fmt.Sprintf("%s:%s", furl, "1984")
	}
	if u.Scheme == "" {
		furl = fmt.Sprintf("http://%s", furl)

	}
	return Dial(furl)
}

// CreateTransaction creates a brand new transaction
func (c *ArweaveClient) CreateTransaction(w *Wallet, amount string, data string, target string) (*Transaction, error) {
	lastTx, err := c.LastTransaction(w.Address())
	if err != nil {
		return nil, err
	}
	price, err := c.GetReward([]byte(data))
	if err != nil {
		return nil, err
	}

	// Non encoded transaction fields
	tx := Transaction{
		lastTx:   lastTx,
		owner:    w.pubKey.N,
		quantity: amount,
		target:   target,
		data:     data,
		reward:   price,
		tags:     make([]map[string]interface{}, 0),
	}

	return &tx, nil
}

// SendTransaction formats the transactions (base64url encodes the necessary fields)
// marshalls the Json and sends it to the arweave network
func (c *ArweaveClient) SendTransaction(tx *Transaction) (string, error) {
	if len(tx.signature) == 0 {
		return "", errors.New("transaction missing signature")
	}
	txn := tx.Format()
	serialized, err := json.Marshal(txn)
	if err != nil {
		return "", err
	}
	return c.Commit(serialized)
}
