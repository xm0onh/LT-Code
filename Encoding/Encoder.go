package Encoding

import (
	//"LT-Code/Decoding"
	"bytes"
	"github.com/bits-and-blooms/bloom"

	//"crypto/sha256"
	C "LT-Code/Cryptography"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/helper/xor"
	//	"github.com/hashicorp/vault/sdk/helper/xor"
	"go.dedis.ch/kyber"

)
/*
func (MicroBlockSlice MicroBlockSliceStruct) EncodeMacroBlock(seq int,degree int) EncDecStructs.Droplet {
//	BufferForEncodingFirstElementofMcBlockInDegree := new(bytes.Buffer)
//	K := len(MicroBlockSlice.MicBlock)
//	fmt.Println("Degree is", K)
	//degree := GetDegree(K,seed)
	//degree:=ChooseDegree(DegreeDistribution)
//	rand.Seed(int64(degree)) // initialize local pseudorandom generator
//	firstElementinDegree:=rand.Intn(len(MicroBlockSlice.MicBlock))
	//firstElementinDegree := degree[0]
//	JstoNEcodingtime:=time.Now()
//	json.NewEncoder(BufferForEncodingFirstElementofMcBlockInDegree).Encode(MicroBlockSlice.MicBlock[firstElementinDegree])
//	fmt.Println("JsonEncoding time for a single MicroBlock", time.Since(JstoNEcodingtime))
	droplet := new(EncDecStructs.Droplet)
	droplet.SeqMicroBlockSlice = make([]bool, len(MicroBlockSlice.MicBlock))
//	droplet.XorMicroBlocks = BufferForEncodingFirstElementofMcBlockInDegree.Bytes()
	droplet.BlockId = MicroBlockSlice.BlockID
	droplet.Seq = seq
	var seen map[int]int
	seen := make(map[int]int, 15)
	if degree == 1 {
		rand.Seed(int64(1)) // initialize local pseudorandom generator
		indxTobeXored:=rand.Intn(len(MicroBlockSlice.MicBlock))
		if _,ok:=seen[indxTobeXored];ok{

		}

	}
	for i:=0;i<degree;i++ {

		rand.Seed(int64(i)) // initialize local pseudorandom generator
		indxMicroBlockTobeIncluded:=rand.Intn(len(MicroBlockSlice.MicBlock))
		fmt.Println("len Microbloclslice.MicBlock is", len(MicroBlockSlice.MicBlock))
		droplet.SeqMicroBlockSlice[indxMicroBlockTobeIncluded] = true
		droplet.HashMicroBlocks = append(droplet.HashMicroBlocks, HashingStructSha256(MicroBlockSlice.MicBlock[indxMicroBlockTobeIncluded]))

			//	McBlkBuffer1,McBlkBuffer2:=JsonByteEncoder(MicroBlockSlice.MicBlock[d],MicroBlockSlice.MicBlock[d+1])
			McBlkBuffer1 := JsonByteEncoder(MicroBlockSlice.MicBlock[indxMicroBlockTobeIncluded])
			XorMicroBuffer := bytes.NewBuffer(droplet.XorMicroBlocks)
			//equivalize  sizes before XoR
			McBlkBuffer1Padded, McBlkBuffer2Padded := CheckAndAddPadding(McBlkBuffer1, XorMicroBuffer)
			var err error
			TimeToXorAMicroblock:=time.Now()
			droplet.XorMicroBlocks, err = xor.XORBytes(McBlkBuffer1Padded.Bytes(), McBlkBuffer2Padded.Bytes())
			if err != nil {
				fmt.Println("Error is", err)
			}
			fmt.Println("Time to XOR a Microblock is", time.Since(TimeToXorAMicroblock))

	}
	return *droplet
}


// A common operation is to XOR entire code blocks together with other blocks.
// When this is done, padding bytes count as 0 (that is XOR identity), and the
// destination block will be modified so that its data is large enough to
// contain the result of the XOR.
/*

*/
//If padding is required returns padded buffer  and true else returns first bytes.buffer with false
func CheckAndAddPadding(reqBodyBytes0, reqBodyBytes1 []byte) ([]byte, []byte) {
	if len(reqBodyBytes0) > len(reqBodyBytes1) {
		padlen := len(reqBodyBytes0) - len(reqBodyBytes1)

		padding := make([]byte, padlen)

		reqBodyBytes1 = append(reqBodyBytes1, padding...)
		return reqBodyBytes0, reqBodyBytes1
	} else {
		padlen := len(reqBodyBytes1) - len(reqBodyBytes0)
		//	fmt.Println("Padlen is", padlen)
		padding := make([]byte, padlen)

		reqBodyBytes0 = append(reqBodyBytes0, padding...)
		return reqBodyBytes0, reqBodyBytes1

	}
	//return reqBodyBytes0, reqBodyBytes1
}

func JsonByteEncoder(MicroBlk1 MicroBlock) *bytes.Buffer {
	reqBodyBytes0 := new(bytes.Buffer)
	//reqBodyBytes1 := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes0).Encode(MicroBlk1)
	if err != nil {
		fmt.Println("json  encoder error is ", err)
	}
	//json.NewEncoder(reqBodyBytes1).Encode(MicroBlk2)
	return reqBodyBytes0
}

/*
func GetDegree(maxNumberOfBlksToBeEncoded, Seed int) []int {
	degreeSlice := make([]int, maxNumberOfBlksToBeEncoded)
	for i := 0; i < maxNumberOfBlksToBeEncoded; i++ {
		degreeSlice[i] = i
	}
	//fmt.Println("Degree slice is", degreeSlice)
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(int64(Seed))
	//r := rand.New(s) // initialize local pseudorandom generator
	degree := rand.Intn(maxNumberOfBlksToBeEncoded)
	if degree == 0 {
		degree = 1
	}
//	fmt.Println("Degree value is", degree)
	rand.Shuffle(len(degreeSlice), func(i, j int) { degreeSlice[i], degreeSlice[j] = degreeSlice[j], degreeSlice[i] })
//	fmt.Println("returning degree slice", degreeSlice[:degree])
	return degreeSlice[:degree]
}

func HashingStructSha256(o interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))
	return h.Sum(nil)
}

*/
func (microblock MicroBlock) GenerateLubyTransformBlock(microblockSlice []MicroBlock, indices []int, priv kyber.Scalar, nodeID string) Droplet {
	//droplet.SeqMicroBlockSlice[indx]
	//fmt.Println("Microblock bytes are",MicroblockinBytes )
	//	droplet.SeqMicroBlockSlice[indx]= true // indices of microblocks

	//	droplet.XorMicroBlocks= make([]byte,len(MicroblockinBytes))
	//	fmt.Println("XorMicroblocks are", droplet.XorMicroBlocks)
	droplet := Initializedroplet(microblockSlice, nodeID)
	//droplet.BlockId = block.BlockId
	//droplet.Seq = block.Seq
	//fmt.Println("Droplet inside is", droplet)
	for _, v := range indices {
		if v < len(microblockSlice) {
			//	fmt.Println("v index is", v)
			droplet.xor(microblockSlice[v])
			droplet.SeqMicroBlockSlice[v] = true
			//	fmt.Println("Droplet inside loop is", droplet)
		}
	}
	//	droplet.xor(microblockSlice[0],0)
	//fmt.Println("Droplet inside GenerateLucy is", droplet)
	hash := C.CalcHash(droplet.XorMicroBlocks)
	droplet.DropletHash = hash
	droplet.Sig = C.SignMsg(droplet.DropletHash, priv)
	droplet.Bloom=bloom.BloomFilter{}
	return droplet
}

func (droplet *Droplet) xor(block MicroBlock) {
	MicroblockinBytes := JsonByteEncoder(block).Bytes()
	//sum := sha256.Sum256(MicroblockinBytes)
	//	hash:=C.CalcHash(droplet.XorMicroBlocks)
	//	droplet.DropletHash = hash
	var inc int
	if len(droplet.XorMicroBlocks) > len(MicroblockinBytes) {
		inc = len(droplet.XorMicroBlocks) - len(MicroblockinBytes)
		MicroblockinBytes = append(MicroblockinBytes, make([]byte, inc)...)

	} else {
		inc = len(MicroblockinBytes) - len(droplet.XorMicroBlocks)
		droplet.XorMicroBlocks = append(droplet.XorMicroBlocks, make([]byte, inc)...)
		//droplet.XorMicroBlocks,_=xor.XORBytes(droplet.XorMicroBlocks, MicroblockinBytes)
	}
	droplet.XorMicroBlocks, _ = xor.XORBytes(droplet.XorMicroBlocks, MicroblockinBytes)
	//	fmt.Println("Droplet inside xor is", droplet)
	//return droplet
}

func Initializedroplet(microblockSlice []MicroBlock, nodeID string) Droplet {
	sliceLen := len(microblockSlice)
	//	C.CalcHash(XorMicroBlocks)

	return Droplet{
		DropletHash:        []byte{},
		SeqMicroBlockSlice: make([]bool, sliceLen, sliceLen),
		//HashMicroBlocks:    make([][]byte, sliceLen, sliceLen),
		XorMicroBlocks: make([]byte, 0),
		BlockId:        microblockSlice[0].BlockId,
		Sig:            nil,
		NodeID:         nodeID,
	}

}

func GenerateBloomFilter(dropletSlice []Droplet, CommitteeSize int)  []Droplet{
	bloom := bloom.NewWithEstimates(uint(CommitteeSize), 0.0000001)
	for _, droplete := range dropletSlice {
		bloom.Add(droplete.DropletHash)
	}

	addBloomFilterToDropletes(dropletSlice, *bloom)
return dropletSlice
}

func addBloomFilterToDropletes(dropletSlice []Droplet, bloom bloom.BloomFilter) {
	for _, droplete := range dropletSlice {
		droplete.Bloom = bloom
	}
}

