package Encoding

import (
	"go.dedis.ch/kyber/v3"
	bn256 "go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/sign/bls"
)

//SeqMicroBlockSlice is a bit array: An on bit will show which Microblock is encoded inside droplet
//HashMicroBlocks Slice of Hash of each included Microblock
//XorMicroBlocks []byte	XORed value of MicroBlocks
//Seq int		Seq of droplet
//BlockId int	Macro Block Id

type Droplet struct {
	DropletHash        []byte
	SeqMicroBlockSlice []bool
	//HashMicroBlocks    [][]byte
	XorMicroBlocks []byte
	Seq            int
	BlockId        int
	Sig            []byte
	NodeID         string
	// Bloom          *bloom.BloomFilter
}

//var DropletSliceMap  map[int][]Droplet

//func InitDropletSliceMap() map[int][]Droplet{
//	DropletSliceMap = make(map[int][]Droplet)
//	return DropletSliceMap

//}

func (D Droplet) Verify(IdTOPbKeyMap map[string]kyber.Point) bool {
	//Hash := ConvertB32toByte(D.DropletHash)
	PubKey := IdTOPbKeyMap[D.NodeID]
	err := bls.Verify(bn256.NewSuite(), PubKey, D.DropletHash, D.Sig)
	if err == nil {
		return true
	}
	return false
}

// Convert [32]byte to []byte
func ConvertB32toByte(B [32]byte) []byte {
	X := B
	return X[:]
}
