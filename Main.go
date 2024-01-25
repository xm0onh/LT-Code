package main

import (
	"github.com/bits-and-blooms/bloom"
	Con "github.com/xm0onh/LT-Code/ConActInterface"
	kzg "github.com/xm0onh/LT-Code/KZG"
	"github.com/xm0onh/LT-Code/Timer"

	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"

	Crypt "github.com/xm0onh/LT-Code/Cryptography"
	"github.com/xm0onh/LT-Code/Decoding"
	"github.com/xm0onh/LT-Code/Encoding"
	"github.com/xm0onh/LT-Code/Net"
)

func main() {
	//
	/////////////////Network Implementation///////////
	gob.Register(Encoding.Droplet{})
	gob.Register(Encoding.Request{})
	gob.Register(Timer.TimerStruct{})
	gob.Register(kzg.KZGStatus{})
	gob.Register(kzg.KZGRequest{})
	gob.Register(kzg.KZGVerify{})
	gob.Register(Net.SerializableKZGVerify{})
	gob.Register(bloom.BloomFilter{})

	////////Key setup and Loading//////////
	regions := []string{"us-east-2", "eu-central-1"}
	fmt.Println("Getting Node IDs")
	NodeIds := Net.GetIDs(regions)
	// Crypt.KeySetup(len(NodeIds))
	fmt.Println("Node Ids for all regions are", NodeIds)
	//	fmt.Println("Getting my ID")
	MyID := Net.GetmyID()
	fmt.Println("My ID is", MyID)
	//fmt.Println("Getting my Idx")
	MyIndx := Decoding.GetIndex(MyID, NodeIds)
	fmt.Println("Getting keys")
	// Crypt.KeySetup(len(NodeIds))
	// Crypt.Load_CommitNum(CommitLen)
	privKey, _ := Crypt.Load_Own_keys("Priv"+strconv.Itoa(MyIndx), "Pub"+strconv.Itoa(MyIndx))
	pubkeys := Crypt.Load_PubKeys(len(NodeIds))
	fmt.Println("Private Key is", privKey)
	fmt.Println("Public Key Slice is", pubkeys)
	conAct := Con.CreateConActor(len(NodeIds), privKey)
	conAct.PubKeySlice = pubkeys
	conAct.PrivateKey = privKey
	/////////Loading IDs/////////
	conAct.MyID = MyID
	// conAct.Peers = Net.Ec2IpExtractor("Role1", "Root-nodes", "us-east-2")
	conAct.RequestorIDs = Net.Ec2IpExtractor("us-east-1", "Role3", "Requestors")
	requestorIps := Net.EC2IPsForAllRegions(regions, "Role3", "Requestors")
	conAct.RequestorIPs = append(conAct.RequestorIPs, requestorIps...)

	//conAct.Peers=append(conAct.Peers,conAct.Primary)
	peers := Net.EC2IPsForAllRegions(regions, "Role1", "Root-nodes")
	conAct.Peers = append(conAct.Peers, peers...)

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
	// conAct.BlockProposalPort="18001"
	// CommitteeSize := (len(conAct.IDs) / 3) + 3
	fmt.Println("--->", os.Args)
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

	const maxResponders = 20                          // Maximum number of responders
	var responderPorts = make([]string, maxResponders) // Map to store responders' IPs and their corresponding ports
	basePort := 18002                                  // Starting port number

	// port2 := "18003"
	myIPAdd, errip := Net.IPaddress()

	//	myAdd:="192.168.4.1"
	if errip != nil {
		fmt.Println(errip)
	}
	for i := 0; i < maxResponders; i++ {
		port := strconv.Itoa(basePort + i)
		responderPorts[i] = port
		Net.InitListener(myIPAdd, port, conAct, CommitteeSize)
	}

	// var ConNetInf ConActAndNetInterface.NetworkToConActInterface
	// Net.InitListener(myIPAdd, conAct.MsgsPort, conAct, CommitteeSize)
	// Net.InitListener(myIPAdd, port2, conAct, CommitteeSize)
	for indx, value := range pubkeys {
		fmt.Println("Index in pubkeySlice is", indx)
		conAct.MapIDToPbKey[conAct.IDs[indx]] = value

	}
	fmt.Println("MapIDToPbKey is", conAct.MapIDToPbKey)
	conAct.Decoder = Decoding.InitDecoder()

	////////////////////////////Decoder//////
	fmt.Println("Requestor IDs are", conAct.RequestorIDs)
	var dropletSlice []Encoding.Droplet
	if !Net.IfIamArequestor(conAct.RequestorIDs, MyID) {

		macroblockSlice := Encoding.GenerateMacroBlocks(numberofMacroBlocks, numberOfMicroBlocks, NumberOfTransactionInEachMicroBlock)
		for _, value := range *macroblockSlice {
			dropletSlice = Encoding.GenerateDropletSlice(value, numberOfMicroBlocks, numberOfMicroBlocks/2, 0.1, conAct.PrivateKey, conAct.MyID)
			fmt.Println("Len Droplet Slice is", len(dropletSlice))
			// fmt.Println("Droplet hash", dropletSlice[0].DropletHash)
			// dropletSlice = Encoding.GenerateBloomFilter(dropletSlice, CommitteeSize)
			fmt.Println("The value of Block ID is", value.BlockID)
			conAct.Decoder.MacroBlockIDToDropletSliceMap[value.BlockID] = dropletSlice

		}
		fmt.Println("Done with Decoding!")

		fmt.Println("Start the KZG Trusted Setup")
		conAct.KZGSetup = *kzg.InitKZG(dropletSlice)
		// fmt.Println("KZG Commitment is", comm)
		// fmt.Println("KZG z is", z)
		// fmt.Println("KZG y is", y)

	}
	/////////////////////KZG Commitment//////////

	/////////////////////Encoder//////////
	// fmt.Println("Requestor IDs are", conAct.RequestorIDs)
	// fmt.Println("My ID is", conAct.MyID)
	//if len(conAct.NodeIdToDialConnMap) == 0 {

	//}
	conAct.TimeCalc = time.Now()
	if Net.IfIamArequestor(conAct.RequestorIDs, conAct.MyID) {

		fmt.Println(" I am a requestor!")
		kzgReq := kzg.CreateKZGRequest()
		totalBlocks := numberofMacroBlocks
		// blocksForNode1 := totalBlocks / 2 // Half of the blocks for Node 1
		// blocksForNode2 := totalBlocks - blocksForNode1 // Remaining blocks for Node 2
		blocksPerResponder := totalBlocks / maxResponders
		for responderIndex := 0; responderIndex < 20; responderIndex++ {
			go func(responderIndex int) {
				startBlock := responderIndex*blocksPerResponder + 1
				endBlock := startBlock + blocksPerResponder - 1
				if responderIndex == maxResponders-1 {
					endBlock = totalBlocks // Ensure the last responder covers all remaining blocks
				}
				request := Encoding.CreateReq(startBlock, endBlock, conAct.MyID, conAct.PrivateKey)
				// request2 := Encoding.CreateReq(blocksForNode1, totalBlocks, conAct.MyID, conAct.PrivateKey)
				// fmt.Println("Request Sig is", request.Sig)
				// fmt.Println("Request Hash is", request.RHash)
				// fmt.Println("conAct.ID is", conAct.IDs)
				port := responderPorts[responderIndex]

				for ID, IP := range conAct.IDToIPMPResponders {
					conAct.NodeIdToDialConnMapResponders[ID] = Net.DialNode(IP, port)
					// conAct.NodeIdToDialConnMapResponders[ID] = Net.DialNode(IP, port2)
					fmt.Println("conAct.NodeIdToDialConnMap[value] is", conAct.NodeIdToDialConnMapResponders[ID])
					fmt.Println("Sending request msg!!!!!!")

				}
				conAct.AddEncodertoNodeIDMap(conAct.NodeIdToDialConnMapResponders)

				for ID, IP := range conAct.IDToIPMPResponders {
					go Net.KZGZSender(conAct.NodeIdToDialConnMapResponders[ID], kzgReq, IP, ID, port, &conAct.NodeIdToDialConnMapResponders, &conAct.NodeIDToEncoderMap)
					fmt.Println(<-conAct.KZGVerficationStatus)
					// fmt.Println("Test after Verification")
					Net.MsgSender(conAct.NodeIdToDialConnMapResponders[ID], request, IP, ID, port, &conAct.NodeIdToDialConnMapResponders, &conAct.NodeIDToEncoderMap)

				}
				// fmt.Println("the request is ", request)
			}(responderIndex)
		}

	}

	idle()

}

func idle() {
	for {
		time.Sleep(1 * time.Minute)
	}

}
