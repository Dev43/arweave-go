package batcher

import (
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/Dev43/arweave-go/chunker"
	"github.com/Dev43/arweave-go/tx"
)

type mockArClient struct {
	LastTx    string
	LastTxErr error

	Tx    func(hash string) *tx.Transaction
	TxErr error
}

func (m *mockArClient) LastTransaction(ctx context.Context, address string) (string, error) {
	return m.LastTx, m.LastTxErr
}
func (m *mockArClient) GetReward(ctx context.Context, data []byte) (string, error) {
	return "1", nil
}
func (m *mockArClient) Commit(ctx context.Context, data []byte) (string, error) {
	return "OK", nil
}
func (m *mockArClient) GetTransaction(ctx context.Context, txID string) (*tx.Transaction, error) {
	return m.Tx(txID), m.TxErr
}

func createNewTestTransaction(id []byte, data string, tags []map[string]interface{}) *tx.Transaction {
	tx := tx.NewTransaction("", nil, "", "", ([]byte(data)), "", tags)
	tx.SetID(id)
	return tx
}

func sendTransaction(hash string) *tx.Transaction {
	switch hash {
	case "0xa":
		return createNewTestTransaction([]byte("0xa"), `{"data": "hi", "position": 0}`, []map[string]interface{}{{AppName: ChunkerInformation{PreviousChunk: "", Position: 0}}})
	case "0xb":
		return createNewTestTransaction([]byte("0xb"), `{"data": "hello", "position": 1}`, []map[string]interface{}{{AppName: ChunkerInformation{PreviousChunk: "0xa", Position: 1}}})
	case "0xc":
		return createNewTestTransaction([]byte("0xc"), `{"data": "there", "position": 2}`, []map[string]interface{}{{AppName: ChunkerInformation{PreviousChunk: "0xb", Position: 2, IsHead: true}}})
	}
	return nil
}

func TestGetAllChunks(t *testing.T) {
	s := strings.NewReader("abc")
	ch, err := chunker.NewChunker(s)
	if err != nil {
		log.Fatal(err)
	}
	ch.SetMaxChunkSize(3)

	cases := []struct {
		mockClient *mockArClient
	}{
		{
			&mockArClient{
				TxErr: nil,
				Tx:    sendTransaction,
			},
		},
	}

	for _, c := range cases {
		bc := NewBatchCombiner(c.mockClient)

		chunks, err := bc.GetAllChunks("0xc")
		if err != nil {
			t.Fatal(err)
		}
		if len(chunks) != 3 {
			t.Fatal(errors.New("invalid chunks lengths"))
		}
	}

}
