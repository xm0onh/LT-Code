package Net

import "C"
import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"reflect"

	E "github.com/xm0onh/LT-Code/Encoding"
	kzg "github.com/xm0onh/LT-Code/KZG"

	//C "github.com/xm0onh/LT-Code/ConActInterface"
	CNI "github.com/xm0onh/LT-Code/ConActAndNetInterface"
)

////Get conactor as an argument

func InitListener(myAdd, BftMsgsPort string, conInterface CNI.NetworkToConActInterface, CommitteeSize int) {
	//port 8003
	serverProposal, err := net.Listen("tcp", myAdd+":"+BftMsgsPort)
	//serverBfTMsgs, err := net.Listen("tcp", myAdd+":"+BftMsgsPort)

	if err != nil {
		fmt.Println(err)
	}
	log.Println("server Proposal", serverProposal)
	log.Println("listening to the port for block", BftMsgsPort)
	log.Println("listening to the port for time duration msgs", BftMsgsPort)
	go SyncListenerLoop(serverProposal, conInterface, CommitteeSize)
	//go TimertListenerLoop(serverBfTMsgs, conInterface,CommitteeSize)

}

func SyncListenerLoop(L net.Listener, conInterface CNI.NetworkToConActInterface, CommitteeSize int) {
	//	defer L.Close()
	fmt.Println("inside BlockListenerLoop")

	for {
		conn, err := L.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go Blockhandleconnection(conn, conInterface, CommitteeSize)

	}

}

/*
func  TimertListenerLoop(L net.Listener,conInterface CNI.NetworkToConActInterface, CommitteeSize int) {
	//	defer L.Close()
	fmt.Println("inside BlockListenerLoop")

	for {
		conn, err := L.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go TimeDurationleconnection(conn, conInterface,CommitteeSize)

	}


}

*/

// func Blockhandleconnection(conn net.Conn, conInterface CNI.NetworkToConActInterface, CommitteeSize int) {

// 	//	fmt.Println("inside Blockhandleconnection decoder")
// 	var RecType1 E.VerifyEntity
// 	var RecType2 kzg.KZGZSender
// 	decoder := gob.NewDecoder(conn)
// 	//////To be done: Read deadline
// 	for {

// 		err := decoder.Decode(&RecType1)
// 		err2 := decoder.Decode(&RecType2)
// 		// fmt.Println("Received Block is", RecType1)
// 		fmt.Println("type is", reflect.TypeOf(RecType2))
// 		if err != nil || err2 != nil {
// 			if err.Error() == "gob: unknown type id or corrupted data" {
// 				fmt.Println("Error during BftMsgshandleconnection", err)

// 			} else {
// 				fmt.Println("here")
// 				conn.Close()
// 				conn = nil
// 				return
// 			}
// 		}
// 		fmt.Println("type is", reflect.TypeOf(RecType1))
// 		fmt.Println("type is", reflect.TypeOf(RecType2))

// 		fmt.Println("Conn in Server.go ", conn)
// 		ipaddress, ErrBin := GetIPaddFromConn(conn)
// 		if !ErrBin {
// 			log.Fatal("Conn corruption in Server")
// 		}
// 		//		fmt.Println("Received Msg", RecType1)
// 		fmt.Println("Committee Size is", CommitteeSize)
// 		fmt.Println("IP address is ", ipaddress)
// 		go conInterface.PassMsgToActor(RecType2, CommitteeSize, ipaddress)

// 		go conInterface.PassMsgToActor(RecType1, CommitteeSize, ipaddress)

// 	}

// }

func Blockhandleconnection(conn net.Conn, conInterface CNI.NetworkToConActInterface, CommitteeSize int) {
	decoder := gob.NewDecoder(conn)
	for {
		// Decode a type identifier string first
		var dataType string
		err := decoder.Decode(&dataType)
		if err != nil {
			fmt.Println("Error decoding data type:", err)
			conn.Close()
			return
		}

		switch dataType {
		case "VerifyEntity":
			var RecType1 E.VerifyEntity
			err = decoder.Decode(&RecType1)
			if err != nil {
				handleDecodeError(err, conn)
				continue // Go to the next loop iteration
			}
			handleReceivedData(RecType1, conn, conInterface, CommitteeSize)

		case "KZGZSender":
			var RecType2 kzg.KZGZSender
			err = decoder.Decode(&RecType2)
			if err != nil {
				handleDecodeError(err, conn)
				continue // Go to the next loop iteration
			}
			handleReceivedData(RecType2, conn, conInterface, CommitteeSize)

		case "KZGZVerifier":
			var RecType2 SerializableKZGVerify
			err = decoder.Decode(&RecType2)

			if err != nil {
				handleDecodeError(err, conn)
				continue // Go to the next loop iteration
			}
			kzgVerifyStruct := RecType2.ToKZGVerify()
			fmt.Println("KZGZVerifier is", kzgVerifyStruct)
			fmt.Println("Type", reflect.TypeOf(kzgVerifyStruct))
			handleReceivedData(kzgVerifyStruct, conn, conInterface, CommitteeSize)
		default:
			fmt.Println("Unknown data type received:", dataType)
		}
	}
}

func handleDecodeError(err error, conn net.Conn) {
	if err.Error() == "gob: unknown type id or corrupted data" {
		fmt.Println("Error during BftMsgshandleconnection:", err)
	} else {
		fmt.Println("Decode error occurred.")
		conn.Close()
	}
}

func handleReceivedData(data interface{}, conn net.Conn, conInterface CNI.NetworkToConActInterface, CommitteeSize int) {
	ipaddress, ErrBin := GetIPaddFromConn(conn)
	if !ErrBin {
		log.Fatal("Conn corruption in Server")
	}

	go conInterface.PassMsgToActor(data, CommitteeSize, ipaddress)
}

/*


func  TimeDurationleconnection(conn net.Conn,conInterface CNI.NetworkToConActInterface, CommitteeSize int) {

	//	fmt.Println("inside Blockhandleconnection decoder")
	var RecType1 Timer.TimerStruct
	for{
		decoder := gob.NewDecoder(conn)

		err := decoder.Decode(&RecType1)
		//fmt.Println("Received Block is", RecType1)
		if err != nil {
			if err.Error() == "gob: unknown type id or corrupted data" {
				fmt.Println("Error during BftMsgshandleconnection", err)

			} else {
				conn.Close()
				conn = nil
				return
			}
		}
		fmt.Println("Conn in Server.go ", conn)
			ipaddress,ErrBin:=GetIPaddFromConn(conn)
		if !ErrBin{
			log.Fatal("Conn corruption in Server")
		}
		//		fmt.Println("Received Msg", RecType1)
		fmt.Println("Committee Size is", CommitteeSize)
		fmt.Println("IP address is ", ipaddress)

		go conInterface.PassMsgToActor(RecType1, CommitteeSize, ipaddress)

	}

}

*/

func GetIPaddFromConn(conn net.Conn) (string, bool) {
	fmt.Println(" conn is ", conn)
	if add, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		return add.IP.String(), true

	}
	return "", false
}
