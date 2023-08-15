package Cryptography

import (
	"crypto/sha256"
	//"go.dedis.ch/kyber/share"
	"go.dedis.ch/kyber"
	bn256 "go.dedis.ch/kyber/pairing/bn256"
	"go.dedis.ch/kyber/sign/bls"
	"strconv"
	//	"CleanThresholdSig/Verification"
	//	D "CleanThresholdSig/Def"
	//	"strconv"
)

func SignMsg(MsgHash []byte, priv kyber.Scalar) []byte {
	Sig, _ := bls.Sign(bn256.NewSuite(), priv, MsgHash)
	return Sig
}

func VerifySignature(MsgHash, Sig []byte, pubkey kyber.Point) bool {
	err := bls.Verify(bn256.NewSuite(), pubkey, MsgHash, Sig)
	if err == nil {
		return true
	}
	return false
}

func ConcatenateByteArr(NodeId string, Seq int) []byte {
	conCat := make([]byte, 0)
	NodeIdinByteArr := []byte(NodeId)
	//	ViewNumByte:=[]byte(strconv.Itoa(viewNum))
	SeqByte := []byte(strconv.Itoa(Seq))
	SliceofByteArray := [][]byte{
		NodeIdinByteArr,
		//	ViewNumByte,
		SeqByte,
	}
	for _, value := range SliceofByteArray {
		conCat = append(conCat, value...)
	}
	return conCat
}

func CalcHash(data []byte) []byte {
	h := sha256.New()
	//str := TransactionEaggStr(TR)
	//h.Write([]byte(str))
	h.Write(data)
	//	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return h.Sum(nil)

}
