package RobustSolitarDistribution

import (
	"math/rand"
	"sort"
)

/*
func ChooseDegree(DegreeDistribution []float32) int{
	UintDegreeDistribution:=floatTouintDegreeDistribution(DegreeDistribution)
	//fmt.Println("UintDegreeDistribution",  UintDegreeDistribution)
	SlicesOfChoices:=createSliceOfChoices(UintDegreeDistribution)
	chooser, _:=wr.NewChooser(SlicesOfChoices...)
	result := chooser.Pick().(int)
	return result

}

func floatTouintDegreeDistribution(DegreeDistribution []float32)  []uint{
	UintDegreeDistribution:=make([]uint,0,len(DegreeDistribution))
	for _,v :=range DegreeDistribution{
		UintDegreeDistribution=append(UintDegreeDistribution,uint(v))
	}
	return UintDegreeDistribution

}

func createSliceOfChoices(floatTouintDegreeDistribution []uint) []wr.Choice{
	SliceChoices:=make([]wr.Choice,len(floatTouintDegreeDistribution),len(floatTouintDegreeDistribution))
	fmt.Println("floatTouintDegreeDistribution is", floatTouintDegreeDistribution)
	for idx, value :=range floatTouintDegreeDistribution{
		SliceChoices[idx].Item=idx
		SliceChoices[idx].Weight=value
	}
	return SliceChoices
}
*/

// pickDegree returns the smallest index i such that cdf[i] > r
// (r a random number from the random generator)
// cdf must be sorted in ascending order.

func PickDegree(random *rand.Rand, cdf []float64) int {
	r := random.Float64()
	d := sort.SearchFloat64s(cdf, r)
	if cdf[d] > r {
		return d
	}

	if d < len(cdf)-1 {
		return d + 1
	} else {
		return len(cdf) - 1
	}
}

// sampleUniform picks num numbers from [0,max) uniformly.
// There will be no duplicates.
// If num >= max, simply returns a slice with all indices from 0 to max-1
// without touching the random number generator.
// The returned slice is sorted.

func SampleUniform(random *rand.Rand, num, max int) []int {
	if num >= max {
		picks := make([]int, max)
		for i := 0; i < max; i++ {
			picks[i] = i
		}
		return picks
	}

	picks := make([]int, num)
	seen := make(map[int]bool)
	for i := 0; i < num; i++ {
		p := random.Intn(max)
		for seen[p] {
			p = random.Intn(max)
		}
		picks[i] = p
		seen[p] = true
	}
	sort.Ints(picks)
	return picks
}
