package ConActInterface

import (
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/xm0onh/LT-Code/Decoding"
	Enc "github.com/xm0onh/LT-Code/Encoding"
	kzg "github.com/xm0onh/LT-Code/KZG"
	N "github.com/xm0onh/LT-Code/Net"
	"github.com/xm0onh/LT-Code/Timer"

	//"errors"
	"net"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"

	//"sync"
	"time"
	//N "github.com/xm0onh/LT-Code/Net"
	//. "github.com/xm0onh/LT-Code/ConActAndNetInterface"
)

//type MsgPassingInterface interface {
//	PassMsgToActor(event interface{})
//}

type TimeAndBlockCountStruct struct {
	Duration   time.Duration
	BlockCount int
}

type ConActor struct {

	//BkKeeper		 *Bookkeeper
	RequestorIPs                  []string
	RequestorIDs                  []string
	Decoder                       Decoding.Decoder
	Suite                         bn256.Suite
	PubKeySlice                   []kyber.Point
	NodeIdToDialConnMapRequestors map[string]net.Conn
	NodeIdToDialConnMapResponders map[string]net.Conn
	//BlockTODropletMap map[int]
	DialBlockConnSlice          []net.Conn
	IDs                         []string
	MapIDToPbKey                map[string]kyber.Point
	Peers                       []string
	Addresses                   []string
	MyID                        string
	PrivateKey                  kyber.Scalar
	MsgsPort                    string
	TimeCalc                    time.Time
	TimeAndBlockCount           TimeAndBlockCountStruct
	ResponderRootNodes          []string
	IDToIPMP                    map[string]string
	IDToIPMPResponders          map[string]string
	IDToIPMPRequesters          map[string]string
	MacroBlocksTimeDurations    []int64
	ResponderTimeDurationSlice  []int64
	RequesterTimeDurationSlice  []int64
	RequestCounter              int64
	RequestResponseTimeCounter  int64
	DropletCounter              int64
	NodeIDToEncoderMap          map[string]*gob.Encoder
	KZGSetup                    kzg.KZGSetup
	BloomFilterVerificationTime int64
}

func CreateConActor(mapSize int, privKey kyber.Scalar) *ConActor {
	return &ConActor{
		//	BkKeeper: 		bk,
		RequestorIPs:                  make([]string, 0, 100),
		RequestorIDs:                  make([]string, 0, 100),
		Decoder:                       Decoding.InitDecoder(),
		Suite:                         *bn256.NewSuite(),
		PubKeySlice:                   make([]kyber.Point, 0),
		NodeIdToDialConnMapRequestors: make(map[string]net.Conn),
		NodeIdToDialConnMapResponders: make(map[string]net.Conn),
		ResponderRootNodes:            make([]string, 0, 100),
		DialBlockConnSlice:            make([]net.Conn, 0, 100),
		IDs:                           make([]string, 0, 100),
		MapIDToPbKey:                  make(map[string]kyber.Point),
		Peers:                         make([]string, 0, 100),
		Addresses:                     make([]string, 0, 10),
		MyID:                          "",
		PrivateKey:                    privKey,
		IDToIPMP:                      make(map[string]string),
		IDToIPMPResponders:            make(map[string]string),
		IDToIPMPRequesters:            make(map[string]string),
		MsgsPort:                      "",
		TimeCalc:                      time.Time{},
		TimeAndBlockCount: struct {
			Duration   time.Duration
			BlockCount int
		}{Duration: 0, BlockCount: 0},
		MacroBlocksTimeDurations:    make([]int64, 0, 100),
		ResponderTimeDurationSlice:  make([]int64, 0, 100),
		RequesterTimeDurationSlice:  make([]int64, 0, 100),
		RequestCounter:              0,
		RequestResponseTimeCounter:  0,
		DropletCounter:              0,
		NodeIDToEncoderMap:          make(map[string]*gob.Encoder),
		BloomFilterVerificationTime: 0,
	}
}

func (c *ConActor) PassMsgToActor(event interface{}, committeeSize int, sourceIp string) {
	fmt.Println("I am Responder and my Node id is:", c.MyID)

	switch event := event.(type) {
	case Enc.Request:
		v := event.Verify(c.MapIDToPbKey)

		//Idx,_:=c.GetMyIndx()
		//var cn NetworkToConActInterface
		MyIndx := Decoding.GetIndex(c.MyID, c.ResponderRootNodes)
		fmt.Println("Responser Root Node ->", c.ResponderRootNodes)
		fmt.Println("event end Blck is", event.EndBlockId)
		fmt.Println("event strt Blck is", event.StartBlockId)
		fmt.Println("NodeIdToDialConnMap is", c.NodeIdToDialConnMapRequestors)
		fmt.Println("Source IP is", sourceIp)
		fmt.Println("v is", v)
		if v {
			c.RequestCounter = c.RequestCounter + 1
			currentTime := time.Now()
			if len(c.NodeIdToDialConnMapRequestors) == 0 {
				for ID, IP := range c.IDToIPMPRequesters {
					conn := N.DialNode(IP, c.MsgsPort)
					c.NodeIdToDialConnMapRequestors[ID] = conn

				}
				c.AddEncodertoNodeIDMap(c.NodeIdToDialConnMapRequestors)

			}

			if event.EndBlockId > event.StartBlockId {
				/* Dimension of MacroBlockIDToDropletSlice
				// rows := len(c.Decoder.MacroBlockIDToDropletSliceMap)
				// fmt.Println("Number of rows:", rows)

				// // Assuming the 2D slice is uniform (every row has the same number of columns)
				// if rows > 0 {
				// 	columns := len(c.Decoder.MacroBlockIDToDropletSliceMap[0])
				// 	fmt.Println("Number of columns:", columns)
				// } else {
				// 	fmt.Println("No columns since there are no rows.")
				// }
				*/
				for i := event.StartBlockId; i <= event.EndBlockId; i++ {
					fmt.Println("Sending Each Droplet with seq i", i)
					/////How many droplets each node has to send (depends on the max value of j)/////////
					fmt.Println("Len MacroblockID to droplete Slice is", i, len(c.Decoder.MacroBlockIDToDropletSliceMap[i]))
					for j := 0; j < 2; j++ {
						fmt.Println("Requester nodeID is", event.NodeId)
						fmt.Println("Requester Connection is", c.NodeIdToDialConnMapRequestors[event.NodeId])
						fmt.Println("The index of the droplete sent is", MyIndx+j*len(c.ResponderRootNodes))

						N.MsgSender(c.NodeIdToDialConnMapRequestors[event.NodeId], c.Decoder.MacroBlockIDToDropletSliceMap[i-1][MyIndx+j*len(c.ResponderRootNodes)], sourceIp, event.NodeId, c.MsgsPort, &c.NodeIdToDialConnMapRequestors, &c.NodeIDToEncoderMap)
						time.Sleep(20 * time.Millisecond)
						conn := N.DialNode(sourceIp, c.MsgsPort)
						c.NodeIdToDialConnMapRequestors[event.NodeId] = conn
					}

					///////////////////////////////
					// go N.MsgSender(c.NodeIdToDialConnMapRequestors[event.NodeId], c.Decoder.MacroBlockIDToDropletSliceMap[i][MyIndx], sourceIp, c.MyID, c.MsgsPort, &c.NodeIdToDialConnMapRequestors)

				}
			} else {
				//	for i:=0;i<len(c.Decoder.MacroBlockIDToDropletSliceMap[event.StartBlockId]);i++{
				for i := 0; i < committeeSize; i++ {
					//go N.MsgSender(conn, Msg, peer, nodeID, port, IdToConnMap)
					//		fmt.Println("Len c.Decoder.MacroBlockIDToDropletSliceMap[event.StartBlockId] is", len(c.Decoder.MacroBlockIDToDropletSliceMap[event.StartBlockId]))
					//		fmt.Println("c.Decoder.MacroBlockIDToDropletSliceMap[event.StartBlockId][i] is", c.Decoder.MacroBlockIDToDropletSliceMap[event.StartBlockId][i])
					time.Sleep(10 * time.Millisecond)

					//go
					N.MsgSender(c.NodeIdToDialConnMapRequestors[event.NodeId], c.Decoder.MacroBlockIDToDropletSliceMap[event.StartBlockId][i], sourceIp, event.NodeId, c.MsgsPort, &c.NodeIdToDialConnMapRequestors, &c.NodeIDToEncoderMap)
					fmt.Println("Decoder sent with index", i)
				}
			}
			// Decoding.SendDroplets(event,c.NodeIdToDialConnMap[event.NodeId],c.Decoder.MacroBlockIDToDropletSliceMap[event.MacroBlkId][MyIndx],sourceIp,c.MyID,c.MsgsPort,&c.NodeIdToDialConnMap)
			// go	N.MsgSender(c.NodeIdToDialConnMap[event.NodeId],c.Decoder.MacroBlockIDToDropletSliceMap[event.MacroBlkId][MyIndx],sourceIp,c.MyID,c.MsgsPort,&c.NodeIdToDialConnMap)
			c.RequestResponseTimeCounter += int64(time.Since(currentTime))
			if c.RequestCounter == int64(len(c.RequestorIDs)) {
				fmt.Println("Total Responder Time is", c.RequestResponseTimeCounter)
				for i := 0; i < len(c.RequestorIDs); i++ {
					//	N.TimeDurationMsgSender(c.NodeIdToDialConnMapRequestors[c.RequestorIDs[i]],c.RequestResponseTimeCounter,sourceIp,c.MyID,c.MsgsPort,&c.NodeIdToDialConnMapRequestors)
					ResponderTimerStruct := Timer.TimerStruct{}
					ResponderTimerStruct.Duration = c.RequestResponseTimeCounter
					ResponderTimerStruct.IsRequesterDuration = false
					ReqUesterIP, ErrBin := N.GetIPaddFromConn(c.NodeIdToDialConnMapRequestors[c.RequestorIDs[i]])

					if !ErrBin {
						log.Fatal("Conn corruption with Requester")
					}
					N.MsgSender(c.NodeIdToDialConnMapRequestors[c.RequestorIDs[i]], ResponderTimerStruct, ReqUesterIP, c.RequestorIDs[i], c.MsgsPort, &c.NodeIdToDialConnMapRequestors, &c.NodeIDToEncoderMap)
				}
			}
		}
	//	event.Verify()
	//	c.addVote(event)
	case Enc.Droplet:

		v := event.Verify(c.MapIDToPbKey)
		// fmt.Println("Received Droplet is", event)
		fmt.Println("Droplet verification is", v)
		// strtTime := time.Now()
		// bloomVerification := event.Bloom.Test(event.DropletHash)
		// endTime := time.Since(strtTime).Nanoseconds()
		// c.BloomFilterVerificationTime = c.BloomFilterVerificationTime + endTime
		// if v && bloomVerification {
		// 	c.Decoder.AddDropletToSlice(committeeSize, event, c.TimeCalc, &c.NodeIdToDialConnMapRequestors, c.RequestorIDs, c.MsgsPort, &c.NodeIDToEncoderMap)
		// }
		c.DropletCounter = c.DropletCounter + 1

		fmt.Println("Droplete countr is", c.DropletCounter)
		if c.DropletCounter == int64(committeeSize) {
			fmt.Println("BloomFilter Verification Time for a single block is", c.BloomFilterVerificationTime)
		}
		//	c.AddAndSendTimer(event)
		//case D.Block:
		//	c.verifyBlockAndSendVote(event)
		//View Change case need to be done
	case kzg.KZGStatus:
		if event.Status {
			fmt.Println("KZG Verification is successful")
		} else {
			fmt.Println("KZG Verification is unsuccessful")
		}

	case kzg.KZGRequest:
		fmt.Println("KZG Request is received")
		fmt.Println(event.Z)
		Y := c.KZGSetup.EvaluatePolynomial(event.Z)
		fmt.Println("Y -->", Y)
		c.KZGSetup.Z = event.Z
		c.KZGSetup.Y = Y
		c.KZGSetup.GenerateProof()

		kzgVerfyStruct := kzg.CreateKZGVerifier(c.KZGSetup.TS, c.KZGSetup.Commitment, *c.KZGSetup.Y, *c.KZGSetup.Z, c.KZGSetup.Proof)
		fmt.Println("Proof ->", kzgVerfyStruct)
		for i := 0; i < len(c.RequestorIDs); i++ {
			ReqUesterIP, ErrBin := N.GetIPaddFromConn(c.NodeIdToDialConnMapRequestors[c.RequestorIDs[i]])
			if !ErrBin {
				log.Fatal("Conn corruption with Requester")
			}
			N.KZGZVerifier(c.NodeIdToDialConnMapRequestors[c.RequestorIDs[i]], kzgVerfyStruct, ReqUesterIP, c.RequestorIDs[i], c.MsgsPort, &c.NodeIdToDialConnMapRequestors, &c.NodeIDToEncoderMap)
		}

	case kzg.KZGVerify:
		fmt.Println("KZG Verify is received")
		v := event.VerifyKZGProof()
		if v {
			fmt.Println("KZG Verification is successful")
		} else {
			fmt.Println("KZG Verification is unsuccessful")
		}
	case Timer.TimerStruct:
		if event.IsRequesterDuration {
			c.CollectRespondersTime(event.Duration)

		} else {
			c.CollectRequestersTime(event.Duration)
		}

	}

}

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

//func (c *ConActor) GetMyIndx() (int,error){
//	for idx,value:=range c.IDs{
//		if value==c.MyID{
//			return idx,nil
//		}
//	}
//	return 0,errors.New("Cannot find my idx")
//}

func (c *ConActor) AddEncodertoNodeIDMap(nodeIdToConnectionMap map[string]net.Conn) {
	//MapNodeIDToEncoder:=make(map[string]*gob.Encoder)
	for nodeId, conn := range nodeIdToConnectionMap {
		c.NodeIDToEncoderMap[nodeId] = gob.NewEncoder(conn)
	}
}
