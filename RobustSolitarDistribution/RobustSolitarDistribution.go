package RobustSolitarDistribution

import (
	"fmt"
	"math"
)

/*
func idealDistribution(N int) []float32{
	probabilities := []float32{0,1/float32(N/4)}
	sumProbabilities:=probabilities[1]
	for k:=2; k< N+1;k++ {

		probabilities=append(probabilities,1 / (float32(k) * (float32(k) - 1)))
		sumProbabilities+=1 / (float32(k) * (float32(k) - 1))
	}
	fmt.Println("Sum probabilities,", sumProbabilities)
	return probabilities

}


func DegreeDistribution(N int, failureProbability float32) []float32{
	ROBUST_FAILURE_PROBABILITY:=0.0001
	M := (N / 2) + 1
	R := int(math.Ceil (float64(N)/float64(M)))
	probabilities := []float32{0,1/float32(M)}
	probabilities[0]=0
	for i:=1;i<int(M);i++{
		probabilities=append(probabilities, 1 / (float32(i) * float32(N)))
	}
	probabilities=append(probabilities,float32(math.Log(float64(R) / float64(ROBUST_FAILURE_PROBABILITY)) / float64(N)))
	fmt.Println("Probabilities slice is", probabilities)
	fmt.Println("M is", M)
	fmt.Println("N is ", N)
	fmt.Println("R is", R)

	for i:=M+1;i<N+1;i++{
		probabilities=append(probabilities,0)
	}
	IdealProbabilities:=idealDistribution(N)
	for i:=0;i<len(IdealProbabilities); i++{

		probabilities[i]+=IdealProbabilities[i]
	}
	//probabilities[2]=probabilities[2]
	fmt.Println("Robust Probabilities are", probabilities)
	fmt.Println("Len Robust Probabilities are", len(probabilities))
	sumProbabilities:=float32(0.0)

	for i:=0;i<len(probabilities);i++{
		probabilities[i]=probabilities[i]*100000000000
	}
	for i:=0;i<len(probabilities);i++{
		sumProbabilities+=probabilities[i]
	}
	fmt.Println("Sum Probabilities is", sumProbabilities)

	return probabilities
}

*/
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
	fmt.Println("CDF is", cdf)
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
