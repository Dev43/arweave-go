package batcher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

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

func createNewTestTransaction(id []byte, data string, tags []tx.Tag) *tx.Transaction {
	tx := tx.NewTransaction("", nil, "", "", ([]byte(data)), "")
	for _, tag := range tags {
		tx.AddTag(tag.Name, tag.Value)
	}
	tx.SetID(id)
	return tx
}

func sendTransaction(hash string) *tx.Transaction {
	tag1, _ := json.Marshal(ChunkInformation{PreviousChunk: "", Position: 0})
	tag2, _ := json.Marshal(ChunkInformation{PreviousChunk: "0xa", Position: 1})
	tag3, _ := json.Marshal(ChunkInformation{PreviousChunk: "0xb", Position: 2, IsHead: true})
	switch hash {
	case "0xa":
		return createNewTestTransaction([]byte("0xa"), `{"data": "hi", "position": 0}`, []tx.Tag{tx.Tag{Name: AppName, Value: string(tag1)}})
	case "0xb":
		return createNewTestTransaction([]byte("0xb"), `{"data": "hello", "position": 1}`, []tx.Tag{tx.Tag{Name: AppName, Value: string(tag2)}})
	case "0xc":
		return createNewTestTransaction([]byte("0xc"), `{"data": "there", "position": 2}`, []tx.Tag{tx.Tag{Name: AppName, Value: string(tag3)}})
	}
	return nil
}

func TestGetAllChunks(t *testing.T) {
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

func TestChunkRecombination(t *testing.T) {
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
		b := bytes.NewBufferString("")
		err = Recombine(chunks, b)
		if err != nil {
			t.Fatal(err)
		}
		if b.String() != "hihellothere" {
			t.Fatal(errors.New("failed at recombining string"))
		}
	}

}
