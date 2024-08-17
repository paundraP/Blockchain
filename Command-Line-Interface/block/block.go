package block

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"cli/common"
	"cli/pow"
)

// Block keeps block headers
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// Deserialize deserializes a block
func Deserialize(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// New creates and returns a Block
func New(data string, prevBlockHash []byte) *Block {
	block := &common.Block{ // Using common.Block
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := pow.New(block) // Use pow.New to create ProofOfWork instance
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return (*Block)(block)
}

// NewGenesis creates and returns the genesis Block
func NewGenesis() *Block {
	return New("Genesis Block", []byte{})
}
