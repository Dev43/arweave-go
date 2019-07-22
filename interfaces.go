// Package arweave defines interfaces for interacting with the Arweave Blockchain.
package arweave

import (
	"context"

	"github.com/Dev43/arweave-go/tx"
	"github.com/Dev43/arweave-go/wallet"
)

type TransactionHandler interface {
	CreateTransaction(w *wallet.Wallet, amount string, data []byte, target string) (*tx.Transaction, error)
	WaitMined(ctx context.Context, tx *tx.Transaction) (*tx.JSONTransaction, error)
	SendTransaction(tx *tx.Transaction) (string, error)
}

type TransactionCreator interface {
	CreateTransaction(w *wallet.Wallet, amount string, data []byte, target string) (*tx.Transaction, error)
}

type TransactionSender interface {
	SendTransaction(tx *tx.Transaction) (string, error)
}

type WalletSigner interface {
	Sign(msg []byte) ([]byte, error)
	Verify(msg []byte, sig []byte) error
}