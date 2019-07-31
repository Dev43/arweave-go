package chunker

import (
	"encoding/json"

	"github.com/Dev43/arweave-go/utils"
)

// EncodedChunk is the data structure representing an encoded chunk on the arweave
type EncodedChunk struct {
	Data     string `json:"data"`
	Position int64  `json:"position"`
}

// EncodedChunkJSON is the intermediary data to encode/decode into JSON
type EncodedChunkJSON struct {
	Data     string `json:"data"`
	Position int64  `json:"position"`
}

// NewEncodedChunkJSON creates a new EncodedChunkJSON struct
func NewEncodedChunkJSON(ec *EncodedChunk) *EncodedChunkJSON {
	return &EncodedChunkJSON{
		Data:     utils.EncodeToBase64([]byte(ec.Data)),
		Position: ec.Position,
	}
}

func (ecj *EncodedChunkJSON) toChunk() (*EncodedChunk, error) {
	decoded, err := utils.DecodeString(ecj.Data)
	if err != nil {
		return nil, err
	}
	return &EncodedChunk{
		Data:     string(decoded),
		Position: ecj.Position,
	}, nil
}

// MarshalJSON marshals as JSON
func (ec *EncodedChunk) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewEncodedChunkJSON(ec))
}

// UnmarshalJSON unmarshals as JSON
func (ec *EncodedChunk) UnmarshalJSON(input []byte) error {
	enc := EncodedChunkJSON{}
	err := json.Unmarshal(input, &enc)
	if err != nil {
		return err
	}
	encoded, err := enc.toChunk()
	if err != nil {
		return err
	}
	*ec = *encoded
	return nil
}
