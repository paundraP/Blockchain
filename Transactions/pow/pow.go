package pow

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"

	"cli/common"
)

const targetBits = 24

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *common.Block
	target *big.Int
}

// New creates and returns a new ProofOfWork instance
func NewPOW(b *common.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// prepareData prepares the block data for hashing
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			[]byte(strconv.FormatInt(pow.block.Timestamp, 16)),
			[]byte(strconv.Itoa(targetBits)),
			[]byte(strconv.Itoa(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs the proof-of-work algorithm
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining a new block")
	for {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

// Validate validates the proof-of-work for the block
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
