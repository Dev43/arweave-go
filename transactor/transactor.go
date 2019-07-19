package transactor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/Dev43/arweave-go/api"
	"github.com/Dev43/arweave-go/tx"
	"github.com/Dev43/arweave-go/wallet"
)

// defaultURL is the local host url
const defaultURL = "http://127.0.0.1:1984"

// defaultPort of the arweave client
const defaultPort = "1984"

// Transactor type, allows one to create transactions
type Transactor struct {
	Client *api.Client
}

// NewTransactor creates a new arweave transactor. You need to pass in a context and a url
// If sending an empty string, the default url is localhosts
func NewTransactor(fullUrl string) (*Transactor, error) {
	if fullUrl == "" {
		c, err := api.Dial(defaultURL)
		if err != nil {
			return nil, err
		}
		return &Transactor{
			Client: c,
		}, nil
	}
	u, err := url.Parse(fullUrl)
	if err != nil {
		return nil, err
	}
	formattedURL := fullUrl
	if u.Port() == "" {
		formattedURL = fmt.Sprintf("%s:%s", fullUrl, defaultPort)
	}
	if u.Scheme == "" {
		formattedURL = fmt.Sprintf("http://%s", formattedURL)
	}
	c, err := api.Dial(formattedURL)

	return &Transactor{
		Client: c,
	}, nil
}

// CreateTransaction creates a brand new transaction
func (tr *Transactor) CreateTransaction(w *wallet.Wallet, amount string, data []byte, target string) (*tx.TransactionBuilder, error) {
	lastTx, err := tr.Client.LastTransaction(w.Address())
	if err != nil {
		return nil, err
	}

	price, err := tr.Client.GetReward([]byte(data))
	if err != nil {
		return nil, err
	}

	// Non encoded transaction fields
	tx := tx.NewTransaction(
		lastTx,
		w.PubKey.N,
		amount,
		target,
		data,
		price,
		make([]map[string]interface{}, 0))

	return tx, nil
}

// SendTransaction formats the transactions (base64url encodes the necessary fields)
// marshalls the Json and sends it to the arweave network
func (tr *Transactor) SendTransaction(tx *tx.TransactionBuilder) (string, error) {
	if len(tx.Signature()) == 0 {
		return "", errors.New("transaction missing signature")
	}
	txn := tx.Format()
	serialized, err := json.Marshal(txn)
	if err != nil {
		return "", err
	}
	return tr.Client.Commit(serialized)
}

func (tr *Transactor) WaitMined(ctx context.Context, tx *tx.TransactionBuilder) (*tx.Transaction, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		receipt, err := tr.Client.GetTransaction(tx.EncodedID())
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
			fmt.Printf("Error retrieving transaction %s \n", err.Error())
		} else {
			fmt.Printf("Transaction not yet mined \n")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
