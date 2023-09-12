package RobustSolitarDistribution

import (
	"math"
)

func RobustSolitonDistribution(n int, m int, delta float64) []float64 {
	pdf := make([]float64, n+1)

	pdf[1] = 3/float64(n) + 1/float64(m)
	total := pdf[1]
	for i := 2; i < len(pdf); i++ {
		pdf[i] = (1 / (float64(i) * float64(i-1)))
		if i < m {
			pdf[i] += 1 / (float64(i) * float64(m))
		}
		if i == m {
			pdf[i] += math.Log(float64(n)/(float64(m)*delta)) / float64(m)
		}
		total += pdf[i]
	}

	cdf := make([]float64, n+1)
	for i := 1; i < len(pdf); i++ {
		pdf[i] /= total
		cdf[i] = cdf[i-1] + pdf[i]
	}
	// fmt.Println("CDF is", cdf)
	return cdf
}

/*
func SolitonDistribution(n int) []float64 {
	cdf := make([]float64, n+1)
	cdf[1] = 1 / float64(n)
	for i := 2; i < len(cdf); i++ {
		cdf[i] = cdf[i-1] + (1 / (float64(i) * float64(i-1)))
	}
	return cdf
}

*/
