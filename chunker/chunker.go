package chunker

import (
	"fmt"
	"io"
	"math"
)

const maxChunkSize = 2
const chunkSizeMB = maxChunkSize << (10 * 2) // 2MB max encoded

type Chunker struct {
	reader       io.ReadSeeker
	totalSize    int64
	totalChunks  int64
	currentChunk int64
	maxChunkSize int64
}

type Chunk struct {
	Data     string `json:"data"`
	Position int64  `json:"position"`
}

func NewChunker(reader io.ReadSeeker) (*Chunker, error) {
	totalSize, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	return &Chunker{
		reader:       reader,
		totalChunks:  calculateTotalChunks(totalSize, chunkSizeMB),
		totalSize:    totalSize,
		currentChunk: 0,
		maxChunkSize: chunkSizeMB,
	}, nil
}

func calculateTotalChunks(totalFileSize, maxChunkSIze int64) int64 {
	return int64(math.Ceil(float64(totalFileSize) / float64(chunkSizeMB)))
}

func (c *Chunker) Size() int64 {
	return c.totalSize
}

func (c *Chunker) SetMaxChunkSize(maxChunkSize int64) {
	c.totalChunks = calculateTotalChunks(c.totalSize, maxChunkSize)
	c.maxChunkSize = maxChunkSize
}

func (c *Chunker) TotalChunks() int64 {
	return c.totalChunks
}

func (c *Chunker) Next() (*Chunk, error) {
	if c.currentChunk >= c.totalChunks {
		return nil, io.EOF
	}

	// offset is i * chunkSizeMB
	offset := c.currentChunk * chunkSizeMB

	currentChunkSize := int64(math.Min(chunkSizeMB, float64(c.totalSize-c.currentChunk*chunkSizeMB)))
	data := make([]byte, currentChunkSize)

	_, err := c.reader.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}
	n, err := c.reader.Read(data)
	if err != nil {
		return nil, err
	}
	if n != int(currentChunkSize) {
		return nil, fmt.Errorf("Did not read the right amount of bytes expected %d actual %d ", currentChunkSize, n)
	}
	chunk := Chunk{
		Data:     string(data),
		Position: c.currentChunk,
	}

	c.currentChunk++

	return &chunk, nil

}

func (c *Chunker) ChunkAll() ([]Chunk, error) {
	c.currentChunk = 0
	chunks := make([]Chunk, c.totalChunks)
	for i := int64(0); i < c.totalChunks; i++ {
		chunk, err := c.Next()
		if err == io.EOF {
			chunks[i] = *chunk
			break
		}
		if err != nil {
			return nil, err
		}
		chunks[i] = *chunk
	}
	return chunks, nil
}

func (c *Chunker) Recombine(chunks []Chunk, w io.WriteSeeker) error {
	if len(chunks) < 1 {
		return fmt.Errorf("no chunks supplied")
	}
	lastChunk := int64(-1)
	offset := int64(0)
	for i := 0; i < len(chunks); i++ {
		currentChunk := chunks[i]
		if currentChunk.Position-lastChunk != 1 {
			return fmt.Errorf("chunks not in order")
		}
		currentChunkSize := len(currentChunk.Data)

		_, err := w.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
		n, err := w.Write([]byte(currentChunk.Data))
		if err != nil {
			return err
		}
		// move the offset
		offset += int64(currentChunkSize)
		if n != int(currentChunkSize) {
			return fmt.Errorf("Did not write the right amount of bytes expected %d actual %d ", currentChunkSize, n)
		}
		lastChunk = currentChunk.Position
	}

	return nil
}
