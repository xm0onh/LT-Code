package kzg

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/arnaucube/kzg-commitments-study"
	k "github.com/arnaucube/kzg-commitments-study"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	Enc "github.com/xm0onh/LT-Code/Encoding"
)

type KZGSetup struct {
	TS         *k.TrustedSetup
	Commitment *bn256.G1
	Coeff      []*big.Int
	Proof      *bn256.G1
	Z          *big.Int
	Y          *big.Int
}

// ConvertHashToFieldElement converts the given hash to a big integer that lies within the finite field defined by R.
func ConvertHashToFieldElement(hash []byte) *big.Int {
	hInt := new(big.Int).SetBytes(hash)
	return hInt.Mod(hInt, k.R)
}

// HashesToPolynomial converts an array of file hashes into a polynomial.
func HashesToPolynomial(hashes [][]byte) []*big.Int {
	coefficients := make([]*big.Int, len(hashes))

	for i, hash := range hashes {
		coefficients[i] = ConvertHashToFieldElement(hash)
	}

	return coefficients
}

func RandomFieldElement() *big.Int {
	modulus := k.R
	z := new(big.Int)
	for {
		b := make([]byte, (modulus.BitLen()+7)/8)
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		z.SetBytes(b)
		if z.Cmp(modulus) < 0 {
			return z
		}
	}
}

func (Setup KZGSetup) EvaluatePolynomial(z *big.Int) *big.Int {
	result := big.NewInt(0)
	for i := len(Setup.Coeff) - 1; i >= 0; i-- {
		result.Mul(result, z)
		result.Add(result, Setup.Coeff[i])
		result.Mod(result, kzg.R)
	}
	return result
}

func InitKZG(dropletSlice []Enc.Droplet) *KZGSetup {

	var hashes [][]byte
	for _, droplet := range dropletSlice {
		hashes = append(hashes, droplet.DropletHash)
	}
	coeff := HashesToPolynomial(hashes)
	ts, err := kzg.NewTrustedSetup(len(coeff))
	if err != nil {
		panic(err)
	}
	c := kzg.Commit(ts, coeff)
	fmt.Println("KZG commitment", c)
	// y := evaluatePolynomial(coeff, z)
	// fmt.Println("Y -->", y)
	// proof, err := kzg.EvaluationProof(ts, coeff, z, y)

	return &KZGSetup{
		TS:         ts,
		Commitment: c,
		Coeff:      coeff,
	}
}

func (K KZGSetup) GenerateProof() {
	proof, err := kzg.EvaluationProof(K.TS, K.Coeff, K.Z, K.Y)
	if err != nil {
		panic(err)
	}
	K.Proof = proof
}
