package chunker

import (
	"fmt"
	"io"
	"math"
)

const KB = 1 << (10 * 1)
const MB = 1 << (10 * 2)
const maxEncodedChunkSize = 750 * KB

type Chunker struct {
	reader              io.Reader
	totalSize           int64
	totalEncodedSize    int64
	totalChunks         int64
	currentChunk        int64
	maxEncodedChunkSize int64
}

type Chunk struct {
	Data     string `json:"data"`
	Position int64  `json:"position"`
}

func NewChunker(reader io.Reader, totalSize int64) (*Chunker, error) {
	totalEncodedSize := getEncodedSize(totalSize)
	return &Chunker{
		reader:              reader,
		totalChunks:         calculateTotalChunks(totalEncodedSize, maxEncodedChunkSize),
		totalSize:           totalSize,
		totalEncodedSize:    totalEncodedSize,
		currentChunk:        0,
		maxEncodedChunkSize: maxEncodedChunkSize,
	}, nil
}

func calculateTotalChunks(totalFileSize, maximumChunkSize int64) int64 {
	return int64(math.Ceil(float64(totalFileSize) / float64(maximumChunkSize)))
}

func (c *Chunker) Size() int64 {
	return c.totalSize
}

func (c *Chunker) EncodedSize() int64 {
	return c.totalEncodedSize
}

func (c *Chunker) SetMaxChunkSize(maxChunkSize int64) {
	c.totalChunks = calculateTotalChunks(c.totalEncodedSize, maxChunkSize)
	c.maxEncodedChunkSize = maxChunkSize
}

func (c *Chunker) TotalChunks() int64 {
	return c.totalChunks
}

func getEncodedSize(size int64) int64 {
	return int64(4 * math.Ceil(float64(size)/3))
}

func getDecodedSize(encSize int64) int64 {
	return int64(math.Ceil(float64(encSize) * 3 / 4))
}

func (c *Chunker) Next() (*Chunk, error) {
	if c.currentChunk >= c.totalChunks {
		return nil, io.EOF
	}

	currentChunkSize := getDecodedSize(int64(math.Min(float64(c.maxEncodedChunkSize), float64(c.totalEncodedSize-c.currentChunk*c.maxEncodedChunkSize))))
	data := make([]byte, currentChunkSize)

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

// Recombine recombines all of the chunks together. It starts with the last chunk (the first one to be created) and works it's way to the first one
func Recombine(chunks []Chunk, w io.Writer) error {
	if len(chunks) < 1 {
		return fmt.Errorf("no chunks supplied")
	}
	lastChunk := chunks[len(chunks)-1].Position
	offset := int64(0)
	for i := len(chunks) - 1; i >= 0; i-- {
		currentChunk := chunks[i]
		if currentChunk.Position-lastChunk > 1 || currentChunk.Position-lastChunk < 0 {
			return fmt.Errorf("chunks not in order")
		}
		currentChunkSize := len(currentChunk.Data)

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
