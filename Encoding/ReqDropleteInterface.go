package Encoding

import (
	"math/big"

	"go.dedis.ch/kyber/v3"
)

type VerifyEntity interface {
	Verify(IdTOPbKeyMap map[string]kyber.Point) bool
}

type KZGZSender interface {
	SendZ() *big.Int
}
