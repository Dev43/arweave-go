package wallet

import (
	"crypto/sha256"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const keyfileName = "arweave-test.json"
const address = "WWMgP35v3BRciaex-sdy-VRfd254M8bqK52v0zRH0Lc"
const modulus = int64(8675517201791120013)
const signature = "Tqq9J7St5qX3OcItXtVdiIzSzbrY2BQSRhmQoYPs4GItaQp_Y2tIqSTxHAQLh6X7MjmLHEQFh4L_5u3JAErT10Z22yFSLOYXwmmCDg8CApkWKuanNTs_YjkTQzvFvPfYy5qcCZipVPcR4unuOORJ0_naGbCxMgHkypgNcdWIepDkjqpqBVWu2VPRCWN4Mw-w7v58kJAKTV8fZnj1n4uuhGmCpd6_WvZFaRAl-LJr-iYUNu6oeZoSuzeH-3Y8k-n6erLq1sbIU6NAbGZheG0KViQt4kpPQPtkRzED3avACIxg224qhQ-elze4BjJRVHZ-SnMQQ-O2TJ90dniD6YhU0KazpDEYfix33Ev0MXrlCY9gHIhQFMwOkTWZtDv4gN76daV0x4J2MQH8q2HK8axLcRTJC-uTpQ7PlhqXyU1VER7MIXX7UX4LMgq1lkoU9PY123wDkstJ-TOjix_iZOMdGQ7GqahV70mZk458kgHYgxVkC8g0PtBE55MGBTojWzFv-hxfNAYXbq4yLzx2akdSMlbtL2LrFPZ19bDHUBLdb9Lq09fIYvQjcjZYi1QyPRFDhnhfBXRkrdLh9WVKpmRpGVyZLbnAlQ1Wkw6zeUZhjIs2mGAh0WbEcVQzVK1_I5cfwoTAr0z8cFP26eJGfXnJ4rd1xCdlszSwWFb112Suc6c"

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
func TestLoadKey(t *testing.T) {
	b := helperLoadBytes(t, keyfileName)
	w := NewWallet()

	w.LoadKey(b)
	ensureCorrectCryptoValues(t, w)

	w.LoadKeyFromFile(filepath.Join("testdata", keyfileName))
	ensureCorrectCryptoValues(t, w)

}

func ensureCorrectCryptoValues(t *testing.T, w *Wallet) {
	assert.Equal(t, address, w.Address(), "Address is not the same")
	assert.Equal(t, modulus, w.PubKeyModulus().Int64(), "Modulus is not the same")
}

func TestSignature(t *testing.T) {
	w := NewWallet()
	w.LoadKeyFromFile(filepath.Join("testdata", keyfileName))

	toSign := []byte("hello")
	msg := sha256.Sum256(toSign)

	sig, err := w.Sign(msg[:])
	if err != nil {
		t.Fatal(err)
	}

	err = w.Verify(msg[:], sig)
	if err != nil {
		t.Fatal(err)
	}

}

func TestGeneration(t *testing.T) {
	w := GenerateWallet()

	toSign := []byte("hello")
	msg := sha256.Sum256(toSign)

	sig, err := w.Sign(msg[:])
	if err != nil {
		t.Fatal(err)
	}

	err = w.Verify(msg[:], sig)
	if err != nil {
		t.Fatal(err)
	}
}