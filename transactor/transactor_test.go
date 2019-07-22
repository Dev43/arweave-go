package transactor_test

import (
	"context"
	"encoding/base64"
	"math/big"
	"testing"

	"github.com/Dev43/arweave-go/transactor"
	"github.com/Dev43/arweave-go/tx"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

type mockCaller struct {
	LastTx string
	Reward string
	Txn    *tx.Transaction
}

func (m *mockCaller) LastTransaction(ctx context.Context, address string) (string, error) {
	return m.LastTx, nil
}

func (m *mockCaller) GetReward(ctx context.Context, data []byte) (string, error) {
	return m.Reward, nil
}

func (m *mockCaller) Commit(ctx context.Context, data []byte) (string, error) {
	return "TESTOK", nil
}

func (m *mockCaller) GetTransaction(ctx context.Context, txID string) (*tx.Transaction, error) {
	return m.Txn, nil
}

type mockWallet struct {
	Signature         []byte
	TestAddress       string
	TestPubKeyModulus *big.Int
}

func (w *mockWallet) Sign(msg []byte) ([]byte, error) {
	return w.Signature, nil
}

func (w *mockWallet) Verify(msg []byte, sig []byte) error {
	return nil
}

func (w *mockWallet) Address() string {
	return w.TestAddress
}

func (w *mockWallet) PubKeyModulus() *big.Int {
	return w.TestPubKeyModulus
}

func encodeToBase64(toEncode []byte) string {
	return base64.RawURLEncoding.EncodeToString(toEncode)
}

func TestCreateTransaction(t *testing.T) {

	cases := []struct {
		caller   *mockCaller
		wallet   *mockWallet
		quantity string
		target   string
		data     []byte
		tag      []map[string]interface{}
	}{
		{
			&mockCaller{
				LastTx: "0xA",
				Reward: "1000",
				Txn:    nil},
			&mockWallet{
				Signature:         nil,
				TestAddress:       "0xB",
				TestPubKeyModulus: big.NewInt(1),
			},
			"1",
			"0xC",
			[]byte("hello"),
			make([]map[string]interface{}, 0),
		},
	}

	for _, c := range cases {
		tr := transactor.Transactor{Client: c.caller}
		tx, err := tr.CreateTransaction(ctx, c.wallet, c.quantity, c.data, c.target)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, c.quantity, tx.Quantity(), "quantity field does not match")
		assert.Equal(t, c.target, tx.Target(), "target field does not match")
		assert.Equal(t, c.tag, tx.Tags(), "tags field does not match")
		assert.Equal(t, c.caller.LastTx, tx.LastTx(), "last tx field does not match")
		assert.Equal(t, c.caller.Reward, tx.Reward(), "reward field does not match")
		assert.Equal(t, encodeToBase64(c.wallet.PubKeyModulus().Bytes()), tx.Owner(), "owner field does not match")
		assert.Equal(t, encodeToBase64(c.wallet.PubKeyModulus().Bytes()), tx.Owner(), "owner field does not match")
	}

}
