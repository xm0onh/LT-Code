package kzg

import (
	"math/big"

	"go.dedis.ch/kyber/v3"
)

type KZGStatus struct {
	Status bool
}

type KZGRequest struct {
	z *big.Int
}

func CreateKZGRequest() *KZGRequest {
	return &KZGRequest{
		z: RandomFieldElement(),
	}
}
func (T KZGStatus) Verify(IdTOPbKeyMap map[string]kyber.Point) bool {
	return T.Status
}

func (K KZGRequest) SendZ() *big.Int {
	return K.z
}
