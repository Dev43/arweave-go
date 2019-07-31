// Package arweave defines interfaces for interacting with the Arweave Blockchain.
package arweave

import "math/big"

// WalletSigner is the interface needed to be able to sign an arweave
type WalletSigner interface {
	Sign(msg []byte) ([]byte, error)
	Verify(msg []byte, sig []byte) error
	Address() string
	PubKeyModulus() *big.Int
}

// BatcherAppName is the application name for the batcher. It is added to transaction tags to retrieve them easily.
const BatcherAppName = "arweave-go-batcher"
