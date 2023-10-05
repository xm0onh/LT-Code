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

func RandomFieldElement(modulus *big.Int) *big.Int {
	z := new(big.Int)
	for {
		b := make([]byte, (modulus.BitLen()+7)/8) // Number of bytes needed to represent the modulus
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

func GenerateKZGProof(dropletSlice []Enc.Droplet) *bn256.G1 {

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
	// bHash := new(big.Int).SetBytes(dropletSlice[0].DropletHash)
	// z := RandomFieldElement(kzg.R)
	z := ConvertHashToFieldElement(dropletSlice[0].DropletHash)
	y := evaluatePolynomial(coeff, z)
	fmt.Println("Y -->", y)
	proof, err := kzg.EvaluationProof(ts, coeff, z, y)
	if err != nil {
		panic(err)
	}

	v := kzg.Verify(ts, c, proof, z, y)
	fmt.Println("debug v", v)
	if v != true {
		panic("NOT VERIFIED")
	}
	return c
	// fmt.Println("Poly is", poly)
}

// func GenerateSampleData(dropletSlice []Enc.Droplet) {
// 	var hashes [][]byte
// 	for _, droplet := range dropletSlice {
// 		hashes = append(hashes, droplet.DropletHash)
// 	}
// }

// func Example() {
// 	file1Content := []byte("file1")
// 	file2Content := []byte("file2")

// 	// Compute their SHA-256 hashes.
// 	hash1 := sha256.Sum256(file1Content)
// 	hash2 := sha256.Sum256(file2Content)
// 	fmt.Println(hash1)
// 	fmt.Println(hash2)
// 	// Convert the hashes to a polynomial.
// 	polynomial := HashesToPolynomial([][]byte{hash1[:], hash2[:]})

// 	// Print the polynomial.
// 	for _, coeff := range polynomial {
// 		// This will print each coefficient of the polynomial.
// 		// In this case, it'll be the representation of the hash values of file1 and file2 respectively.
// 		println(coeff.String())
// 	}
// }
