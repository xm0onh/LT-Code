package kzg

import (
	"math/big"

	kzg "github.com/arnaucube/kzg-commitments-study"
)

// ConvertHashToFieldElement converts the given hash to a big integer that lies within the finite field defined by R.
func ConvertHashToFieldElement(hash []byte) *big.Int {
	hInt := new(big.Int).SetBytes(hash)
	return hInt.Mod(hInt, kzg.R)
}

// HashesToPolynomial converts an array of file hashes into a polynomial.
func HashesToPolynomial(hashes [][]byte) []*big.Int {
	coefficients := make([]*big.Int, len(hashes))

	for i, hash := range hashes {
		coefficients[i] = ConvertHashToFieldElement(hash)
	}

	return coefficients
}

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
