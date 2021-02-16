package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/Dev43/arweave-go/utils"
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
	pubKey    *rsa.PublicKey
}

// NewWallet creates a new wallet instance
func NewWallet() *Wallet {
	return &Wallet{}
}

// GenerateWallet generates a new JWK wallet.
func GenerateWallet() *Wallet {
	reader := rand.Reader
	rsaKey, _ := rsa.GenerateKey(reader, 4096)
	w := &Wallet{}

	w.key = &gojwk.Key{
		Kty: "RSA",
		N: utils.EncodeToBase64(rsaKey.N.Bytes()),
		E: utils.EncodeToBase64(big.NewInt(int64(rsaKey.E)).Bytes()),
		D: utils.EncodeToBase64(rsaKey.D.Bytes()),
	}
	w.pubKey = rsaKey.Public().(*rsa.PublicKey)
	// Take the "n", in bytes and hash it using SHA256
	h := sha256.New()
	h.Write(rsaKey.N.Bytes())

	// Finally base64url encode it to have the resulting address
	w.address = utils.EncodeToBase64(h.Sum(nil))
	w.publicKey = utils.EncodeToBase64(rsaKey.N.Bytes())
	return w
}

// Address returns the address of the account
func (w *Wallet) Address() string {
	return w.address
}

// PubKeyModulus returns the modulus of the RSA public key
func (w *Wallet) PubKeyModulus() *big.Int {
	return w.pubKey.N
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

// LoadKeyFromFile loads and Arweave RSA key from a file to our wallet
func (w *Wallet) LoadKeyFromFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return w.LoadKey(b)
}

// LoadKey loads an Arweave RSA key into our wallet
func (w *Wallet) LoadKey(rsaKeyBytes []byte) error {

	key, err := gojwk.Unmarshal(rsaKeyBytes)
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
	w.pubKey = pubKey
	// Take the "n", in bytes and hash it using SHA256
	h := sha256.New()
	h.Write(pubKey.N.Bytes())

	// Finally base64url encode it to have the resulting address
	w.address = utils.EncodeToBase64(h.Sum(nil))
	w.publicKey = utils.EncodeToBase64(pubKey.N.Bytes())
	w.key = key

	return nil
}
