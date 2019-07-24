package batcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Dev43/arweave-go/chunker"
	"github.com/Dev43/arweave-go/transactor"
	"github.com/Dev43/arweave-go/wallet"
)

const chunkerVersion = "0.0.1"
const AppName = "chunker"

type BatchMaker struct {
	ar        *transactor.Transactor
	wallet    *wallet.Wallet
	reader    io.Reader
	totalSize int64
}

type ChunkerInformation struct {
	PreviousChunk string `json:"previous_chunk"`
	IsHead        bool   `json:"is_head"`
	Version       string `json:"version"`
	Position      int64  `json:"position"`
}

func NewBatch(ar *transactor.Transactor, w *wallet.Wallet, reader io.Reader, totalSize int64) *BatchMaker {

	return &BatchMaker{
		ar:        ar,
		wallet:    w,
		reader:    reader,
		totalSize: totalSize,
	}
}

func (b *BatchMaker) SendBatchTransaction() ([]string, error) {
	txList := []string{}
	ch, err := chunker.NewChunker(b.reader, b.totalSize)
	if err != nil {
		return nil, err
	}
	for i := int64(0); i < ch.TotalChunks(); i++ {
		chunk, err := ch.Next()
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(chunk)
		if err != nil {
			return nil, err
		}

		txBuilder, err := b.ar.CreateTransaction(context.TODO(), b.wallet, "", data, "")
		if err != nil {
			return nil, err
		}
		previousChunk := ""
		if len(txList) > 0 {
			previousChunk = txList[len(txList)-1]
		}
		tags := txBuilder.Tags()
		isHead := false
		if i+1 == ch.TotalChunks() {
			isHead = true
		}
		chunkerInfo := ChunkerInformation{
			PreviousChunk: previousChunk,
			IsHead:        isHead,
			Version:       chunkerVersion,
			Position:      chunk.Position,
		}
		tags = append(tags, map[string]interface{}{AppName: chunkerInfo})
		txBuilder.SetTags(tags)
		tx, err := txBuilder.Sign(b.wallet)
		if err != nil {
			return nil, err
		}
		// hash, err := b.ar.SendTransaction(context.TODO(), tx)
		// if err != nil {
		// 	return err
		// }
		// minedTx, err := b.ar.WaitMined(context.TODO(), tx)
		txList = append(txList, tx.Hash())
	}
	fmt.Printf("Successfully sent batch transactions with head transaction %s and list of transactions: \n - %s \n", txList[len(txList)-1], strings.Join(txList, "\n - "))

	return txList, nil
}
