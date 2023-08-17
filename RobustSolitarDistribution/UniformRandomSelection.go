package RobustSolitarDistribution

import (
	"math/rand"
	"sort"
)



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
