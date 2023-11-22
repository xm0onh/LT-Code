package Decoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/xm0onh/LT-Code/Encoding"
	N "github.com/xm0onh/LT-Code/Net"
	"github.com/xm0onh/LT-Code/Timer"

	//	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/sdk/helper/xor"

	"net"
	"time"
)

func (decoder *Decoder) AddDropletToSlice(committeeSize int, droplet Encoding.Droplet, startTime time.Time, NodeIdToDialConnMapRequestors *map[string]net.Conn, NodesSlice []string, MsgsPort string, IdToEncoderMap *map[string]*gob.Encoder) {
	//decoder.LockMacroBlockIDToDropletSliceMap.Lock()
	if _, ok := decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId]; !ok {
		dropletSlice := make([]Encoding.Droplet, 0, committeeSize)
		decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId] = append(dropletSlice, droplet)
		//	decoder.DecodedMicroBlockcounter[droplet.BlockId] = 0

	} else {
		decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId] = append(decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId], droplet)
	}
	//	decoder.LockMacroBlockIDToDropletSliceMap.Unlock()

	fmt.Println("Block to droplete ID slice length is", len(decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId]), "Block ID is", droplet.BlockId, "Committee Size is", committeeSize)
	//if len(decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId]) >= 2*committeeSize {
	if len(decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId]) >= committeeSize {

		//	decoder.IfDropletIsReadyToBeDecoded[droplet.BlockId]=true
		for i := 0; i < len(decoder.MacroBlockIDToDropletSliceMap[droplet.BlockId]); i++ {
			decoder.Peel(droplet.BlockId)
			if len(decoder.Blockchain.MapMacroBlockNumToMapMiroBlockHashToMicroBlock[droplet.BlockId]) == committeeSize/2 {
				fmt.Println("Number of decoded Microblocks are", len(decoder.Blockchain.MapMacroBlockNumToMapMiroBlockHashToMicroBlock[droplet.BlockId]))
				fmt.Println("Decoded Microblcks are", decoder.Blockchain.MapMacroBlockNumToMapMiroBlockHashToMicroBlock[droplet.BlockId])
				totalTimeTaken := time.Since(startTime)
				RequesterTimeStruct := Timer.TimerStruct{}
				RequesterTimeStruct.Duration = totalTimeTaken.Nanoseconds()
				///Msg sender needs to be adjusted.
				for j := 0; j < len(*NodeIdToDialConnMapRequestors); j++ {
					N.MsgSender((*NodeIdToDialConnMapRequestors)[NodesSlice[j]], RequesterTimeStruct, (*NodeIdToDialConnMapRequestors)[NodesSlice[j]].RemoteAddr().String(), NodesSlice[j], MsgsPort, NodeIdToDialConnMapRequestors, IdToEncoderMap)

				}

				fmt.Println("TotalTime is", totalTimeTaken)
			}
		}

	}
}

func (decoder *Decoder) GetDroplet(MacroBlockID, idx int) *Encoding.Droplet {
	return &decoder.MacroBlockIDToDropletSliceMap[MacroBlockID][idx]

}

func (decoder *Decoder) GetSingleton(MacroBlockId int) (int, int, bool) {
	//counter :=0
	EdgeToBlockinSeq := 0
	//val:=false
	fmt.Println("Inside Get Singleton is")
	//fmt.Println("decoder.MacroBlockIDToDropletSliceMap[MacroBlockId]", decoder.MacroBlockIDToDropletSliceMap[MacroBlockId])
	for idx, droplet := range decoder.MacroBlockIDToDropletSliceMap[MacroBlockId] {
		counter := 0
		//fmt.Println("This droplet inside GetSingleton is", droplet)
		for idxEdgeToBlockinSeq, val := range droplet.SeqMicroBlockSlice {
			//	fmt.Println("Bool value is",val)
			if val == true {
				counter++
				EdgeToBlockinSeq = idxEdgeToBlockinSeq
			}
		}
		if counter == 1 {
			//	fmt.Println("This counter is 1111111111111111")
			//	fmt.Println("This counter is 1111111111111111")
			//MicroBlockHash:=string(droplet.HashMicroBlocks[idx])
			microblock := JsonDecoder(droplet.XorMicroBlocks)
			//fmt.Println("Decoded Microblock is", microblock)
			///////////New Map
			if _, ok := decoder.Blockchain.MapMacroBlockNumToMapMiroBlockHashToMicroBlock[droplet.BlockId]; !ok {
				decoder.Blockchain.MapMacroBlockNumToMapMiroBlockHashToMicroBlock[droplet.BlockId] = make(map[string]Encoding.MicroBlock)

			}
			//////

			//		if _, ok := decoder.Blockchain.MapBlockHashToMicroBlock[string(microblock.MicroBlockHash)]; ok {
			//		fmt.Println("microblock has already present")
			//	}
			decoder.Blockchain.MapMacroBlockNumToMapMiroBlockHashToMicroBlock[droplet.BlockId][string(microblock.MicroBlockHash)] = microblock
			//decoder.Blockchain.MapBlockHashToMicroBlock[string(microblock.MicroBlockHash)] = microblock //Commented because of the line above
			decoder.DecodedMicroBlockcounter[MacroBlockId] += 1
			//decoder.DecodedBlock[BlockId][EdgeToBlockinSeq]=true
			decoder.RemoveNeighborsFromSingleton(idx, EdgeToBlockinSeq, MacroBlockId)

			return idx, EdgeToBlockinSeq, true
		}
	}
	return 0, 0, false
}

//	func  removeEdge(idx int, droplet Encoding.Droplet){
//		droplet.SeqMicroBlockSlice[idx]=false
//	}
//
// /In peel first we get a singleton and add its MicroBlock to the mainChain. Then we find all droplets that has edges to
// /to the respective MicroBlock of singleton. We simply perform xor between those Microblocks and the singleton and remove
// /those edges from droplets.
func (decoder *Decoder) Peel(blockId int) bool {
	_, _, success := decoder.GetSingleton(blockId)
	//fmt.Println("idx,seqIdx,bol", SingleTonIndexinDropletSliceMap,idxBlockWithinSingleton,success)

	if !success {
		fmt.Println("Failed to decode download more")
		return false
	}

	//for idx,droplet:=range decoder.MacroBlockIDToDropletSliceMap[blockId]{
	//	if idx != SingleTonIndexinDropletSliceMap{
	//		if droplet.SeqMicroBlockSlice[idxBlockWithinSingleton]==true{

	//		decoder.MacroBlockIDToDropletSliceMap[blockId][idx].XorMicroBlocks,_ = xor.XORBytes(decoder.MacroBlockIDToDropletSliceMap[blockId][idx].XorMicroBlocks,decoder.MacroBlockIDToDropletSliceMap[blockId][SingleTonIndexinDropletSliceMap].XorMicroBlocks)
	//		removeEdge(idxBlockWithinSingleton,decoder.MacroBlockIDToDropletSliceMap[blockId][idx])
	//		decoder.MacroBlockIDToDropletSliceMap[blockId]=append(decoder.MacroBlockIDToDropletSliceMap[blockId][:idxBlockWithinSingleton],decoder.MacroBlockIDToDropletSliceMap[blockId][idxBlockWithinSingleton+1:]...)
	//		fmt.Println("Droplet SliceMap BlockID is", decoder.MacroBlockIDToDropletSliceMap[blockId][idx].Seq)
	//		decoder.Peel(blockId)
	// 		}
	return true

	//	}
	//}
	//return false

}

func (decoder *Decoder) RemoveNeighborsFromSingleton(SingleTonIndexinDropletSliceMap, idxBlockWithinSingleton, blockId int) {
	singleton := decoder.MacroBlockIDToDropletSliceMap[blockId][SingleTonIndexinDropletSliceMap]

	decoder.MacroBlockIDToDropletSliceMap[blockId][SingleTonIndexinDropletSliceMap] = decoder.MacroBlockIDToDropletSliceMap[blockId][len(decoder.MacroBlockIDToDropletSliceMap[blockId])-1]
	decoder.MacroBlockIDToDropletSliceMap[blockId] = decoder.MacroBlockIDToDropletSliceMap[blockId][:len(decoder.MacroBlockIDToDropletSliceMap[blockId])-1]
	//decoder.MacroBlockIDToDropletSliceMap[blockId]=append(decoder.MacroBlockIDToDropletSliceMap[blockId][:SingleTonIndexinDropletSliceMap],decoder.MacroBlockIDToDropletSliceMap[blockId][SingleTonIndexinDropletSliceMap+1:]...)
	var err error
	for idx, _ := range decoder.MacroBlockIDToDropletSliceMap[blockId] {
		if decoder.MacroBlockIDToDropletSliceMap[blockId][idx].SeqMicroBlockSlice[idxBlockWithinSingleton] == true {
			neighboringXorBlocks, singletonBlock := Encoding.CheckAndAddPadding(decoder.MacroBlockIDToDropletSliceMap[blockId][idx].XorMicroBlocks, singleton.XorMicroBlocks)
			decoder.MacroBlockIDToDropletSliceMap[blockId][idx].XorMicroBlocks, err = xor.XORBytes(neighboringXorBlocks, singletonBlock)
			fmt.Println("Error is", err)
		}
		decoder.MacroBlockIDToDropletSliceMap[blockId][idx].SeqMicroBlockSlice[idxBlockWithinSingleton] = false

	}
}

//func (decoder Decoder) CheckIfBlockForDropleteIsNotDecoded(HashMicroBlock string) bool {
///	if _, ok := decoder.Blockchain.MapBlockHashToMicroBlock[HashMicroBlock]; !ok {
//	return false
//}
//return true
//}

func (decoder *Decoder) GetNumberofDecodedMicroBlockForMacroBlock(MacroBlockId int) int {
	return decoder.DecodedMicroBlockcounter[MacroBlockId]
}

func JsonDecoder(DataInBytes []byte) Encoding.MicroBlock {
	Mc1 := new(Encoding.MicroBlock)
	reader := bytes.NewReader(DataInBytes)
	err := json.NewDecoder(reader).Decode(Mc1)
	if err != nil {
		fmt.Println("Error during Decding is ", err)
	}
	return *Mc1

}

func SendDroplets(request Encoding.Request, conn net.Conn, Msg Encoding.VerifyEntity, peer, nodeID, port string, IdToConnMap *map[string]net.Conn) {
	if request.EndBlockId > request.StartBlockId {
		for i := request.StartBlockId; i <= request.EndBlockId; i++ {
			//go	MsgSender(c.NodeIdToDialConnMap[request.NodeId],c.Decoder.MacroBlockIDToDropletSliceMap[i][MyIndx],sourceIp,c.MyID,c.MsgsPort,&c.NodeIdToDialConnMap)
			//	go N.MsgSender(conn, Msg, peer, nodeID, port, IdToConnMap)

		}
	} else {
		//	go	N.MsgSender(c.NodeIdToDialConnMap[request.NodeId],c.Decoder.MacroBlockIDToDropletSliceMap[request.StartBlockId][MyIndx],sourceIp,c.MyID,c.MsgsPort,&c.NodeIdToDialConnMap)
		//for i:=0;i<36;i++{
		//	go N.MsgSender(conn, Msg, peer, nodeID, port, IdToConnMap)
		//go	MsgSender(c.NodeIdToDialConnMap[request.NodeId],c.Decoder.MacroBlockIDToDropletSliceMap[request.StartBlockId][i],sourceIp,c.MyID,c.MsgsPort,&c.NodeIdToDialConnMap)
		//	}
	}
}

func GetIndex(ID string, IDSlice []string) int {
	ind := 0
	for index, nodeId := range IDSlice {
		if ID == nodeId {
			return index
		}
	}
	return ind
}
