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

func evaluatePolynomial(p []*big.Int, x *big.Int) *big.Int {
	result := big.NewInt(0)
	for i := len(p) - 1; i >= 0; i-- {
		result.Mul(result, x)
		result.Add(result, p[i])
		result.Mod(result, kzg.R)
	}
	return result
}

func InitKZG(dropletSlice []Enc.Droplet) (*k.TrustedSetup, *bn256.G1, *big.Int, *big.Int, *bn256.G1, error) {

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
	z := RandomFieldElement()
	// z := ConvertHashToFieldElement(dropletSlice[0].DropletHash)
	y := evaluatePolynomial(coeff, z)
	fmt.Println("Y -->", y)
	proof, err := kzg.EvaluationProof(ts, coeff, z, y)

	return ts, c, z, y, proof, err
}

func VerifyKZGProof(TS *k.TrustedSetup, commit *bn256.G1, proof *bn256.G1, z, y *big.Int) bool {
	return kzg.Verify(TS, commit, proof, z, y)
}
