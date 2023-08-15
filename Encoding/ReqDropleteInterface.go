package Encoding
import "go.dedis.ch/kyber"

type VerifyEntity interface {
	Verify(IdTOPbKeyMap map[string]kyber.Point)bool
}
