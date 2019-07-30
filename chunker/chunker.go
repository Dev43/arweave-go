package chunker

import (
	"fmt"
	"io"
	"math"
)

const kb = 1 << (10 * 1)
const mb = 1 << (10 * 2)

const maxEncodedChunkSize = 500 * kb

// Chunker struct
type Chunker struct {
	reader       io.Reader
	totalSize    int64
	totalChunks  int64
	currentChunk int64
}

// EncodedChunk is the data structure representing an encoded chunk on the arweave
type EncodedChunk struct {
	Data     string `json:"data"`
	Position int64  `json:"position"`
}

// NewChunker creates a new chunker struct
func NewChunker(reader io.Reader, totalSize int64) (*Chunker, error) {
	return &Chunker{
		reader:       reader,
		totalChunks:  calculateTotalChunks(totalSize, maxEncodedChunkSize),
		totalSize:    totalSize,
		currentChunk: 0,
	}, nil
}

func calculateTotalChunks(totalFileSize, chunkSize int64) int64 {
	return int64(math.Ceil(float64(totalFileSize) / float64(chunkSize)))
}

// Size retrieves the total size of the chunk
func (c *Chunker) Size() int64 {
	return c.totalSize
}

// EncodedSize calculates the Base64RawURL encoded size of our chunk
func (c *Chunker) EncodedSize() int64 {
	return getEncodedSize(c.totalSize)
}

// SetChunkSize sets the chunk size
func (c *Chunker) SetChunkSize(size int64) {
	c.totalChunks = calculateTotalChunks(c.totalSize, size)
}

// TotalChunks returns the total chunks
func (c *Chunker) TotalChunks() int64 {
	return c.totalChunks
}

// Next retrieves the next chunk from the io.Reader
func (c *Chunker) Next() (*EncodedChunk, error) {
	if c.currentChunk >= c.totalChunks {
		return nil, io.EOF
	}
	size := int64(math.Min(float64(maxEncodedChunkSize), float64(c.totalSize-c.currentChunk*maxEncodedChunkSize)))
	data := make([]byte, size)

	n, err := c.reader.Read(data)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, fmt.Errorf("did not read the right amount of bytes expected %d actual %d ", size, n)
	}
	chunk := EncodedChunk{
		Data:     string(data),
		Position: c.currentChunk,
	}

	c.currentChunk++

	return &chunk, nil

}

// ChunkAll is a commodity function designed to returns all the chunks from an io.Reader
func (c *Chunker) ChunkAll() ([]EncodedChunk, error) {
	c.currentChunk = 0
	chunks := make([]EncodedChunk, c.totalChunks)
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

// Recombine recombines all of the chunks. It starts with the last chunk (which is the first one to be created) and works it's way to the first one
func Recombine(chunks []EncodedChunk, w io.Writer) error {
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

func paddingSize(inputSize int64) int64 {
	if inputSize%3 == 0 {
		return 3 - (inputSize % 3)
	}
	return 0
}

func getEncodedSize(size int64) int64 {
	return int64(((4 * size / 3) + 3) & ^3)
}

func getDecodedSize(encSize int64) int64 {
	return int64(math.Ceil(float64(encSize) * 3 / 4))
}
