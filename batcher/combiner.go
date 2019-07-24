package batcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Dev43/arweave-go/chunker"
	"github.com/Dev43/arweave-go/transactor"
)

type BatchCombiner struct {
	arClient transactor.ClientCaller
}

func NewBatchCombiner(client transactor.ClientCaller) *BatchCombiner {
	return &BatchCombiner{
		arClient: client,
	}
}

func (bc *BatchCombiner) GetAllChunks(headChunkAddress string) ([]chunker.Chunk, error) {
	headTx, err := bc.arClient.GetTransaction(context.TODO(), headChunkAddress)
	if err != nil {
		return nil, err
	}
	chunk := chunker.Chunk{}
	tags := headTx.Tags()
	chunkInfo, err := getChunkInfoFromTag(tags)
	if err != nil {
		return nil, err
	}
	if !chunkInfo.IsHead {
		return nil, fmt.Errorf("transaction is not the head chunk transaction")
	}

	err = json.Unmarshal((headTx.RawData()), &chunk)
	if err != nil {
		return nil, err
	}
	chunks := []chunker.Chunk{}
	chunks = append(chunks, chunk)
	return bc.getChunk(chunkInfo.PreviousChunk, chunks)
}

func getChunkInfoFromTag(tags []map[string]interface{}) (*ChunkerInformation, error) {
	for _, tag := range tags {
		ch, ok := tag[AppName]
		if ok {
			chunkInfo, ok := ch.(ChunkerInformation)
			if !ok {
				return nil, fmt.Errorf("could not cast tags to ChunkerInformation")
			}

			return &chunkInfo, nil
		}
	}
	return nil, fmt.Errorf("necessary tag not present in transaction")
}

func (bc *BatchCombiner) getChunk(address string, chunks []chunker.Chunk) ([]chunker.Chunk, error) {
	if address == "" {
		return chunks, nil
	}
	tx, err := bc.arClient.GetTransaction(context.TODO(), address)
	if err != nil {
		return nil, err
	}
	chunk := chunker.Chunk{}
	err = json.Unmarshal([]byte(tx.RawData()), &chunk)
	if err != nil {
		return nil, err
	}
	chunkInfo, err := getChunkInfoFromTag(tx.Tags())
	if err != nil {
		return nil, err
	}
	chunks = append(chunks, chunk)

	return bc.getChunk(chunkInfo.PreviousChunk, chunks)
}

func Recombine(chunks []chunker.Chunk, w io.Writer) error {
	return chunker.Recombine(chunks, w)
}
