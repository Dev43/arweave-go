package arweave

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/ioutil"

	"gopkg.in/square/go-jose.v2"
)

type Wallet struct {
	address   string
	publicKey string
	key       jose.JSONWebKey
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) Public() string {
	return w.publicKey
}

func (w *Wallet) Sign(msg []byte) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS256, Key: w.key}, nil)
	if err != nil {
		return "", err
	}
	object, err := signer.Sign(msg)
	// var payload = []byte("Lorem ipsum dolor sit amet")
	// object, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}
	serialized, err := object.CompactSerialize()
	if err != nil {
		return "", err
	}
	return serialized, nil
}

func (w *Wallet) Verify(msg []byte, sig string) error {

	// Parse the serialized, protected JWS object. An error would indicate that
	// the given input did not represent a valid message.
	// Needs to be serialized
	object, err := jose.ParseSigned(sig)
	if err != nil {
		return err
	}

	// Now we can verify the signature on the payload. An error here would
	// indicate the the message failed to verify, e.g. because the signature was
	// broken or the message was tampered with.
	_, err = object.Verify(w.key.Public().Key)
	if err != nil {
		return err
	}
	return nil

}

// Assumes a normal unencrypted JSK
func (w *Wallet) ExtractKey(fileName string) error {

	key := jose.JSONWebKey{}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = key.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	w.key = key

	// Type cast the key to get the "n"
	rsa, ok := key.Public().Key.(*rsa.PublicKey)
	if !ok {
		return errors.New("could not typecast key to rsa public key")
	}

	// Take the "n", in bytes and hash it using SHA256
	h := sha256.New()
	h.Write(rsa.N.Bytes())

	// Finally base64url encode it to have the resulting address
	w.address = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	w.publicKey = base64.RawURLEncoding.EncodeToString(rsa.N.Bytes())
	return nil
}

func NewWallet() *Wallet {
	return &Wallet{}
}
