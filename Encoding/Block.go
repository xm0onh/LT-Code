package Encoding

import (
	"fmt"

	"go.dedis.ch/kyber/v3"

	//	"time"
	rnd "math/rand"

	RSD "github.com/xm0onh/LT-Code/RobustSolitarDistribution"
)

func GenerateDropletSlice(macroblock MicroBlockSliceStruct, NumberOfMicroBlocks, m int, delta float64, priv kyber.Scalar, nodeID string) []Droplet {
	//n=33, m=15, delta=0.1
	DegreeSliceFloat := RSD.RobustSolitonDistribution(NumberOfMicroBlocks, m, delta)
	fmt.Println("Degree Slice length is floati is", len(DegreeSliceFloat))
	//DegreeSliceFloat:=Rsd.SolitonDistribution(3000)
	degreeSlice := make([]int, 0, len(DegreeSliceFloat))
	for i := 0; i < 3*len(DegreeSliceFloat); i++ {
		degree := RSD.PickDegree(rnd.New(RSD.NewMersenneTwister(int64(120+i))), DegreeSliceFloat)
		degreeSlice = append(degreeSlice, degree)
	}
	uniformlySelectedIndices := make([]int, 0, 3*len(DegreeSliceFloat))
	//droplet:=Encoding.Initializedroplet(macroblock.MicBlock)
	dropletSlice := make([]Droplet, 0, len(degreeSlice))
	//	StarttimeToEncode := time.Now()
	for indx, value := range degreeSlice {
		//	for i:=0;i<value;i++{
		uniformlySelectedIndices = RSD.SampleUniform(rnd.New(rnd.NewSource(int64(100+indx))), value, len(DegreeSliceFloat))
		//}
		//fmt.Println("Macroblock is", macroblock)
		droplet := macroblock.MicBlock[uniformlySelectedIndices[0]].GenerateLubyTransformBlock(macroblock.MicBlock, uniformlySelectedIndices, priv, nodeID)
		dropletSlice = append(dropletSlice, droplet)
	}
	return dropletSlice
}
