package Encoding

import (
	"fmt"

	Crypto "github.com/xm0onh/LT-Code/Cryptography"
	"go.dedis.ch/kyber/v3"
	bn256 "go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/sign/bls"
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
	fmt.Println("Sender Node ID during verification is:", R.NodeId)
	fmt.Println("R.NodeId is", R.NodeId)
	fmt.Println("Sender PubKey is", PubKey)
	fmt.Println("R.Hash is  iHash", R.RHash)
	fmt.Println("R.Sig is", R.Sig)

	err := bls.Verify(bn256.NewSuite(), PubKey, R.RHash, R.Sig)
	fmt.Println("Verification err ->", err)
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
