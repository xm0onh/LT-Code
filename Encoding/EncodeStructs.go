package Encoding

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
)

type MicroBlockSliceStruct struct {
	MicBlock []MicroBlock
	BlockID  int
}

const lettersAndNumbers = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}
	return string(b)
}

func GenerateMacroBlockStruct(numberOfMicroblocks, blockId, NumberOfTransactions int) MicroBlockSliceStruct {
	MacroBlock := new(MicroBlockSliceStruct)
	MicroBlockList := make([]MicroBlock, 0, 50)
	for i := 0; i < numberOfMicroblocks; i++ {
		trList := transactionGenerator(NumberOfTransactions)
		MicroBlockList = append(MicroBlockList, generateMicroBlock(trList, i, blockId))
	}
	MacroBlock.MicBlock = MicroBlockList
	MacroBlock.BlockID = blockId
	return *MacroBlock
}

func GenerateMacroBlocks(numberOfMacroBlocks, numberOfMicroblocks, NumberOfTransactions int) *[]MicroBlockSliceStruct {
	SliceMicroBlockSliceStruct := make([]MicroBlockSliceStruct, 0, numberOfMacroBlocks)
	for i := 0; i < numberOfMacroBlocks; i++ {
		SliceMicroBlockSliceStruct = append(SliceMicroBlockSliceStruct, GenerateMacroBlockStruct(numberOfMicroblocks, i, NumberOfTransactions))
	}
	return &SliceMicroBlockSliceStruct
}

func generateTransaction(seq int) Transaction {
	randStringBytes(35)
	randStringBytes(35)
	rand.Intn(100000000000000)
	return Transaction{
		//	From:   "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
		//To:     "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
		From:   randStringBytes(35),
		To:     randStringBytes(35),
		Amount: 1,
		Hash:   randStringBytes(35),
		Seq:    seq,
	}
}

func transactionGenerator(numberOfTransactions int) []Transaction {
	transactionList := make([]Transaction, 0, 0)
	for i := 0; i < numberOfTransactions; i++ {
		transactionList = append(transactionList, generateTransaction(i))
	}
	return transactionList
}

func generateMicroBlock(payload []Transaction, seq, blkId int) MicroBlock {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(payload)
	if err != nil {
		fmt.Println("Error during encoding is ", err)
	}
	h := sha256.New()
	h.Write(reqBodyBytes.Bytes())
	PayloadHash := h.Sum(nil)
	return MicroBlock{
		MicroBlockHash: PayloadHash,
		Payload:        payload,
		Seq:            seq,
		BlockId:        blkId,
	}
}

type Transaction struct {
	From   string
	To     string
	Amount int
	Hash   string
	Seq    int
}

type MicroBlock struct {
	MicroBlockHash []byte
	Payload        []Transaction
	Seq            int
	BlockId        int
}
