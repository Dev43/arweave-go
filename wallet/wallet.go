package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/mendsley/gojwk"
)

var opts = &rsa.PSSOptions{
	SaltLength: rsa.PSSSaltLengthAuto,
	Hash:       crypto.SHA256,
}

// Wallet struct
type Wallet struct {
	address   string
	key       *gojwk.Key
	publicKey string
	PubKey    *rsa.PublicKey
}

// Address returns the address of the account
func (w *Wallet) Address() string {
	return w.address
}

// Sign signs a message using the RSA-PSS scheme with an MGF SHA256 masking function
func (w *Wallet) Sign(msg []byte) ([]byte, error) {
	priv, err := w.key.DecodePrivateKey()
	if err != nil {
		return nil, err
	}
	rng := rand.Reader
	privRsa, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("could not typecast key to %T", rsa.PrivateKey{})
	}

	sig, err := rsa.SignPSS(rng, privRsa, crypto.SHA256, msg, opts)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// Verify verifies the signature for the specific message
func (w *Wallet) Verify(msg []byte, sig []byte) error {
	pub, err := w.key.DecodePublicKey()
	if err != nil {
		return err
	}
	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("could not typecast key to %T", rsa.PublicKey{})
	}

	err = rsa.VerifyPSS(pubKey, crypto.SHA256, msg, sig, opts)
	if err != nil {
		return err
	}
	return nil
}

// ExtractKey extracts the necessary information from the arweave key file.
// It assumes the arweave key file is unencrypted.
func (w *Wallet) ExtractKey(fileName string) error {

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	key, err := gojwk.Unmarshal(b)
	if err != nil {
		return err
	}

	publicKey, err := key.DecodePublicKey()
	if err != nil {
		return err
	}
	pubKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("could not typecast key to %T", rsa.PublicKey{})
	}
	w.PubKey = pubKey
	// Take the "n", in bytes and hash it using SHA256
	h := sha256.New()
	h.Write(pubKey.N.Bytes())

	// Finally base64url encode it to have the resulting address
	w.address = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	w.publicKey = base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())
	w.key = key

	return nil
}

// NewWallet creates a new wallet instance
func NewWallet() *Wallet {
	return &Wallet{}
}
