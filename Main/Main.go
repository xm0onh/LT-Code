package main

import (
	Con "LT-Code/ConActInterface"
	"LT-Code/Timer"
	"github.com/bits-and-blooms/bloom"

	Crypt "LT-Code/Cryptography"
	"LT-Code/Decoding"
	"LT-Code/Encoding"
	"LT-Code/Net"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	
	/////////////////Network Implementation///////////
	gob.Register(Encoding.Droplet{})
	gob.Register(Encoding.Request{})
	gob.Register(Timer.TimerStruct{})
	gob.Register(bloom.BloomFilter{})

	////////Key setup and Loading//////////
	regions := []string{"us-east-1", "eu-central-1", "sa-east-1"}
	fmt.Println("Getting Node IDs")
	NodeIds := Net.GetIDs(regions)
	//	Crypt.KeySetup(len(NodeIds))
	fmt.Println("Node Ids for all regions are", NodeIds)
	//	fmt.Println("Getting my ID")
	MyID := Net.GetmyID()
	fmt.Println("My ID is", MyID)
	//fmt.Println("Getting my Idx")
	MyIndx := Decoding.GetIndex(MyID, NodeIds)
	fmt.Println("Getting keys")
	//Crypt.KeySetup(len(NodeIds))
	//Crypt.Load_CommitNum(CommitLen)
	privKey, _ := Crypt.Load_Own_keys("Priv"+strconv.Itoa(MyIndx), "Pub"+strconv.Itoa(MyIndx))
	pubkeys := Crypt.Load_PubKeys(len(NodeIds))
	fmt.Println("Private Key is", privKey)
	fmt.Println("Public Key Slice is", pubkeys)
	conAct := Con.CreateConActor(len(NodeIds), privKey)
	conAct.PubKeySlice = pubkeys
	conAct.PrivateKey = privKey
	/////////Loading IDs/////////
	conAct.MyID = MyID
	//conAct.Peers= Net.Ec2IpExtractor("Role1", "Root-nodes")
	//	conAct.RequestorIDs = Net.Ec2IpExtractor("us-east-1", "Role3", "Requestors")
	requestorIps := Net.EC2IPsForAllRegions(regions, "Role3", "Requestors")
	for _, requestor := range requestorIps {
		conAct.RequestorIPs = append(conAct.RequestorIPs, requestor)
	}

	//conAct.Peers=append(conAct.Peers,conAct.Primary)
	peers := Net.EC2IPsForAllRegions(regions, "Role1", "Root-nodes")
	for _, peer := range peers {
		conAct.Peers = append(conAct.Peers, peer)
	}
	///////////////////////IDSANDIPSForRespondersANDRequestors////////
	fmt.Println("Just before IDSANDIPSForRespondersANDRequestors!!!!!!!!!!!!!!!")
	ResponderIDSAndIPS := Net.EC2IPsAndIDSForAllRegions(regions, "Role1", "Root-nodes")
	fmt.Println("ResponderIDSAndIPS are", ResponderIDSAndIPS)
	for i := 0; i < len(ResponderIDSAndIPS); i += 2 {
		conAct.IDToIPMPResponders[ResponderIDSAndIPS[i]] = ResponderIDSAndIPS[i+1]
		//	IDToIPMP[IpAndID[i]]=IpAndID[i+1]
	}

	RequestorIDSAndIPS := Net.EC2IPsAndIDSForAllRegions(regions, "Role3", "Requestors")
	fmt.Println("RequestorIDSAndIPS are", RequestorIDSAndIPS)

	for i := 0; i < len(RequestorIDSAndIPS); i += 2 {
		conAct.IDToIPMPRequesters[RequestorIDSAndIPS[i]] = RequestorIDSAndIPS[i+1]
		//	IDToIPMP[IpAndID[i]]=IpAndID[i+1]
	}

	/////////////////////////////////////

	RequestorIDs := Net.Ec2RequestorIDExtractorForAllRegions(regions, "Role3", "Requestors")
	for _, requestorids := range RequestorIDs {
		conAct.RequestorIDs = append(conAct.RequestorIDs, requestorids)
	}
	//conAct.Peers=append(conAct.Peers,conAct.Primary)
	conAct.IDs = NodeIds

	ResponderIDs := Net.Ec2RequestorIDExtractorForAllRegions(regions, "Role1", "Root-nodes")
	for _, responderids := range ResponderIDs {
		conAct.ResponderRootNodes = append(conAct.ResponderRootNodes, responderids)
	}

	///////////// Network Init//////////////////////
	//conAct.BlockProposalPort="18001"
	//CommitteeSize := (len(conAct.IDs) / 3) + 3
	numberofMacroBlocks, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	numberOfMicroBlocks, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	NumberOfTransactionInEachMicroBlock, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	CommitteeSize := 2 * numberOfMicroBlocks
	conAct.MsgsPort = "18002"

	myIPAdd, errip := Net.IPaddress()

	//	myAdd:="192.168.4.1"
	if errip != nil {
		fmt.Println(errip)
	}
	//var ConNetInf ConActAndNetInterface.NetworkToConActInterface
	Net.InitListener(myIPAdd, conAct.MsgsPort, conAct, CommitteeSize)
	for indx, value := range pubkeys {
		fmt.Println("Index in pubkeySlice is", indx)
		conAct.MapIDToPbKey[conAct.IDs[indx]] = value
	}
	//	for idx, peer := range conAct.Peers {

	//		conAct.NodeIdToDialConnMap[conAct.IDs[idx]] = Net.DialNode(peer, conAct.MsgsPort)
	//	}

	//	for idx, requestors := range conAct.RequestorIPs {

	//		conAct.NodeIdToDialConnMap[conAct.IDs[idx]] = Net.DialNode(requestors, conAct.MsgsPort)
	//	}
	conAct.Decoder = Decoding.InitDecoder()
	//	IpAndID:=Net.Ec2IDandIPExtractor("us-east-1")
	//	fmt.Println("IpAndID is ", IpAndID)
	//	for i:=0;i<len(IpAndID);i+=2{
	//		conAct.IDToIPMP[IpAndID[i]]=IpAndID[i+1]
	//	}
	//	fmt.Println("conAct.IDToIPMP is ", conAct.IDToIPMP)
	//	for k,v :=range conAct.IDToIPMP{
	//		conAct.NodeIdToDialConnMap[k] = Net.DialNode(v, conAct.MsgsPort)

	//	}
	//fmt.Println("Nodeto Dial ConnMap in the main is", conAct.NodeIdToDialConnMap)
	////////////////////////////Decoder//////

	if !Net.IfIamArequestor(conAct.RequestorIDs, MyID) {

		//	var testMapMicrobLockHashtoMicroblock map[string]int
		//	testMapMicrobLockHashtoMicroblock = make(map[string]int)
		//macroblock := Encoding.GenerateMacroBlockStruct(numberOfMicroBlocks, 1, NumberOfTransactionInEachMicroBlock)
		macroblockSlice := Encoding.GenerateMacroBlocks(numberofMacroBlocks, numberOfMicroBlocks, NumberOfTransactionInEachMicroBlock)
		for _, value := range *macroblockSlice {
			dropletSlice := Encoding.GenerateDropletSlice(value, numberOfMicroBlocks, numberOfMicroBlocks/2, 0.1, conAct.PrivateKey, conAct.MyID)
			fmt.Println("Len Droplet Slice is", len(dropletSlice))
			dropletSlice=Encoding.GenerateBloomFilter(dropletSlice, CommitteeSize)

			conAct.Decoder.MacroBlockIDToDropletSliceMap[value.BlockID] = dropletSlice

		}
		//for Key,_:= range conAct.Decoder.MacroBlockIDToDropletSliceMap{

		//	}

	}
	//StartTimeToDecode := time.Now()
	//	for i := 0; i < 33; i++ {
	//		decoder.Peel(1)

	//	}
	//EndTimeToDecode := time.Since(StartTimeToDecode)
	//	fmt.Println("Microblocks are, ",decoder.Blockchain.MapBlockHashToMicroBlock)
	//	fmt.Println("Microblocks are, ", decoder.Blockchain.MapBlockHashToMicroBlock)
	//	fmt.Println("Microblocks are, ", len(decoder.Blockchain.MapBlockHashToMicroBlock))

	/////////////////////Encoder//////////
	fmt.Println("Requestor IDs are", conAct.RequestorIDs)
	fmt.Println("My ID is", conAct.MyID)
	//if len(conAct.NodeIdToDialConnMap) == 0 {

	//}
	conAct.TimeCalc = time.Now()
	if Net.IfIamArequestor(conAct.RequestorIDs, conAct.MyID) {

		fmt.Println(" I am a requestor!")
		request := Encoding.CreateReq(1, 2, conAct.MyID, conAct.PrivateKey)

		fmt.Println("Request Sig is", request.Sig)
		fmt.Println("Request Hash is", request.RHash)
		//fmt.Println("NodeIdToDialConnMap is ", conAct.NodeIdToDialConnMapResponders)
		fmt.Println("conAct.ID is", conAct.IDs)
		//	time.Sleep(20 * time.Second)

		for ID, IP := range conAct.IDToIPMPResponders {
			//	if !Net.IfIamArequestor(conAct.RequestorIDs, myIPAdd) {
			///A better option may be to have NodeIdToPeerIPMap instead of extracting IP from Conn.
			//	PeerIP, bolErr := Net.GetIPaddFromConn(conAct.NodeIdToDialConnMap[value])
			//	if !bolErr {
			//		log.Fatal("error while extracting IP from conn.")
			//	}
			//	for k,v :=range conAct.IDToIPMPResponders{
			conAct.NodeIdToDialConnMapResponders[ID] = Net.DialNode(IP, conAct.MsgsPort)
			fmt.Println("conAct.NodeIdToDialConnMap[value] is", conAct.NodeIdToDialConnMapResponders[ID])

			//		}
			fmt.Println("Sending request msg!!!!!!")
			//	}
		}
		conAct.AddEncodertoNodeIDMap(conAct.NodeIdToDialConnMapResponders)

		for ID, IP := range conAct.IDToIPMPResponders {

			Net.MsgSender(conAct.NodeIdToDialConnMapResponders[ID], request, IP, ID, conAct.MsgsPort, &conAct.NodeIdToDialConnMapResponders, &conAct.NodeIDToEncoderMap)
		}
	}

	idle()

}

func idle() {
	for {
		time.Sleep(1 * time.Minute)
	}

}
