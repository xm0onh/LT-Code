package kzg

import "go.dedis.ch/kyber/v3"

type KZGStruct struct {
	Status bool
}

func (T KZGStruct) Verify(IdTOPbKeyMap map[string]kyber.Point) bool {
	return T.Status
}
