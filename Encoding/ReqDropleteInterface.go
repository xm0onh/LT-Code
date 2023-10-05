package Encoding

import (
	"go.dedis.ch/kyber/v3"
)

type VerifyEntity interface {
	Verify(IdTOPbKeyMap map[string]kyber.Point) bool
}
