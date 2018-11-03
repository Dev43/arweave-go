package arweave

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/mendsley/gojwk"
)

type Wallet struct {
	address   string
	key       *gojwk.Key
	publicKey string
}

func (w *Wallet) Address() string {
	return w.address
}

var opts = &rsa.PSSOptions{
	SaltLength: rsa.PSSSaltLengthEqualsHash,
	Hash:       crypto.SHA256,
}

func (w *Wallet) Sign(msg []byte) ([]byte, error) {
	priv, err := w.key.DecodePrivateKey()
	if err != nil {
		return nil, err
	}
	rng := rand.Reader
	privRsa, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("could not typecast key to rsa public key")
	}

	sig, err := rsa.SignPSS(rng, privRsa, crypto.SHA256, msg, opts)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (w *Wallet) Verify(msg []byte, sig []byte) error {
	pub, err := w.key.DecodePublicKey()
	if err != nil {
		return err
	}
	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		fmt.Println("could not typecast to rsa.PublicKey")
	}

	err = rsa.VerifyPSS(pubKey, crypto.SHA256, msg, sig, opts)
	if err != nil {
		return err
	}
	return nil
}

// Assumes a normal unencrypted JSK
func (w *Wallet) Extract(fileName string) error {

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
	rsa, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("could not typecast key to rsa public key")
	}
	// Take the "n", in bytes and hash it using SHA256
	h := sha256.New()
	h.Write(rsa.N.Bytes())

	// Finally base64url encode it to have the resulting address
	w.address = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	w.publicKey = rsa.N.String()
	w.key = key

	return nil
}

func NewWallet() *Wallet {
	return &Wallet{}
}
