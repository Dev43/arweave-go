// Package arweave defines interfaces for interacting with the Arweave Blockchain.
package arweave

// WalletSigner is the interface needed to be able to sign an arweave transaction
type WalletSigner interface {
	Sign(msg []byte) ([]byte, error)
	Verify(msg []byte, sig []byte) error
}
