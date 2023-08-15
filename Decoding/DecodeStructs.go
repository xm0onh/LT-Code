package Decoding

import (
	"LT-Code/Encoding"
	"sync"
)

type Decoder struct {
	//	Map of MacroBlock Id to slice/array of its droplets
	MacroBlockIDToDropletSliceMap     map[int][]Encoding.Droplet
	LockMacroBlockIDToDropletSliceMap sync.RWMutex
	//Index is associated with respective BlockId
	//IfDropletIsReadyToBeDecoded map[int]bool
	//DecodedBlock map[int]map[int]bool
	Blockchain                   ChainStatus
	LockBlockchain               sync.RWMutex
	DecodedMicroBlockcounter     map[int]int
	LockDecodedMicroBlockcounter sync.RWMutex
}

//const mapSize = 1000

type ChainStatus struct {
	//isMainChain bool
	RWLockChain sync.RWMutex
	//	MapBlockHashToMicroBlock map[string]Encoding.MicroBlock
	MapMacroBlockNumToMapMiroBlockHashToMicroBlock map[int]map[string]Encoding.MicroBlock
}

func InitChainStatus() ChainStatus {
	return ChainStatus{
		MapMacroBlockNumToMapMiroBlockHashToMicroBlock: make(map[int]map[string]Encoding.MicroBlock),
		//	MapBlockHashToMicroBlock: make(map[string]Encoding.MicroBlock, mapSize),
	}
}

func InitDecoder() Decoder {
	return Decoder{
		MacroBlockIDToDropletSliceMap: make(map[int][]Encoding.Droplet),
		//	IfDropletIsReadyToBeDecoded: make([]bool,0),
		Blockchain:               InitChainStatus(),
		DecodedMicroBlockcounter: make(map[int]int),
	}
}
