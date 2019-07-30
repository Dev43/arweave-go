package batcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Dev43/arweave-go/chunker"
	"github.com/Dev43/arweave-go/transactor"
	"github.com/Dev43/arweave-go/tx"
)

type BatchCombiner struct {
	arClient transactor.ClientCaller
}

func NewBatchCombiner(client transactor.ClientCaller) *BatchCombiner {
	return &BatchCombiner{
		arClient: client,
	}
}

func (bc *BatchCombiner) GetAllChunks(headChunkAddress string) ([]chunker.EncodedChunk, error) {
	headTx, err := bc.arClient.GetTransaction(context.TODO(), headChunkAddress)
	if err != nil {
		return nil, err
	}
	chunk := chunker.EncodedChunk{}
	tags, err := headTx.Tags()
	if err != nil {
		return nil, err
	}
	chunkInfo, err := getChunkInfoFromTag(tags)
	if err != nil {
		return nil, err
	}
	if !chunkInfo.IsHead {
		return nil, fmt.Errorf("transaction is not the head chunk transaction")
	}

	err = json.Unmarshal(headTx.RawData(), &chunk)
	if err != nil {
		return nil, err
	}
	chunks := []chunker.EncodedChunk{}
	chunks = append(chunks, chunk)
	return bc.getChunk(chunkInfo.PreviousChunk, chunks)
}

func getChunkInfoFromTag(tags []tx.Tag) (*ChunkInformation, error) {
	for _, tag := range tags {
		if tag.Name == AppName {
			chunkInfo := ChunkInformation{}
			err := json.Unmarshal([]byte(tag.Value), &chunkInfo)
			if err != nil {
				return nil, err
			}
			return &chunkInfo, nil
		}
	}
	return nil, fmt.Errorf("necessary tag not present in transaction")
}

func (bc *BatchCombiner) getChunk(address string, chunks []chunker.EncodedChunk) ([]chunker.EncodedChunk, error) {
	if address == "" {
		return chunks, nil
	}
	tx, err := bc.arClient.GetTransaction(context.TODO(), address)
	if err != nil {
		return nil, err
	}
	chunk := chunker.EncodedChunk{}
	err = json.Unmarshal([]byte(tx.RawData()), &chunk)
	if err != nil {
		return nil, err
	}
	tags, err := tx.Tags()
	if err != nil {
		return nil, err
	}
	chunkInfo, err := getChunkInfoFromTag(tags)
	if err != nil {
		return nil, err
	}
	chunks = append(chunks, chunk)

	return bc.getChunk(chunkInfo.PreviousChunk, chunks)
}

// Recombine recombines the chunks and writes it to an io.Writer
func Recombine(chunks []chunker.EncodedChunk, w io.Writer) error {
	return chunker.Recombine(chunks, w)
}
