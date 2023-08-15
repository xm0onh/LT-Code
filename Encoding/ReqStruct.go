package Encoding

import (
	Crypto "LT-Code/Cryptography"
	"fmt"
	"go.dedis.ch/kyber"
	bn256 "go.dedis.ch/kyber/pairing/bn256"
	"go.dedis.ch/kyber/sign/bls"
)

type Request struct {
	StartBlockId int
	EndBlockId   int
	MacroBlkId   int
	NodeId       string
	RHash        []byte
	Sig          []byte
}

func (R Request) Verify(IdTOPbKeyMap map[string]kyber.Point) bool {
	PubKey := IdTOPbKeyMap[R.NodeId]
	fmt.Println("R.NodeId is", R.NodeId)
	fmt.Println("Sender PubKey is", PubKey)
	fmt.Println("R.Hash is  iHash", R.RHash)
	fmt.Println("R.Sig is", R.Sig)

	err := bls.Verify(bn256.NewSuite(), PubKey, R.RHash, R.Sig)
	if err == nil {
		return true
	}
	return false
}

func CreateReq(StartBlkId, EndBlkId int, nodeID string, privKey kyber.Scalar) *Request {
	hash := Crypto.ConcatenateByteArr(nodeID, StartBlkId)
	sig := Crypto.SignMsg(hash, privKey)
	return &Request{
		StartBlockId: StartBlkId,
		EndBlockId:   EndBlkId,
		NodeId:       nodeID,
		RHash:        hash,
		Sig:          sig,
	}

}
