package kzg

import (
	"math/big"

	"github.com/arnaucube/kzg-commitments-study"
	k "github.com/arnaucube/kzg-commitments-study"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type KZGStatus struct {
	Status bool
}

type KZGRequest struct {
	Z *big.Int
}

type KZGVerify struct {
	TS         k.TrustedSetup
	Commitment bn256.G1
	Y          big.Int
	Z          big.Int
	Proof      bn256.G1
}

func CreateKZGVerifier(TS k.TrustedSetup, Commitment bn256.G1, Y big.Int, Z big.Int, Proof bn256.G1) *KZGVerify {
	return &KZGVerify{
		TS:         TS,
		Commitment: Commitment,
		Y:          Y,
		Z:          Z,
		Proof:      Proof,
	}
}
func CreateKZGRequest() *KZGRequest {
	return &KZGRequest{
		Z: RandomFieldElement(),
	}
}

func (K KZGRequest) SendZ() *big.Int {
	return K.Z
}

func (KVerify KZGVerify) VerifyKZGProof() bool {
	return kzg.Verify(&KVerify.TS, &KVerify.Commitment, &KVerify.Proof, &KVerify.Z, &KVerify.Y)
}
