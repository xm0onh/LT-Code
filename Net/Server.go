package Net

import "C"
import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"reflect"

	E "github.com/xm0onh/LT-Code/Encoding"

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

func Blockhandleconnection(conn net.Conn, conInterface CNI.NetworkToConActInterface, CommitteeSize int) {

	//	fmt.Println("inside Blockhandleconnection decoder")
	var RecType1 E.VerifyEntity
	fmt.Println("Checking for the error about droplet types", reflect.TypeOf(RecType1))
	decoder := gob.NewDecoder(conn)
	//////To be done: Read deadline
	for {

		err := decoder.Decode(&RecType1)
		// fmt.Println("Received Block is", RecType1)
		if err != nil {
			if err.Error() == "gob: unknown type id or corrupted data" {
				fmt.Println("Error during BftMsgshandleconnection", err)

			} else {
				conn.Close()
				conn = nil
				return
			}
		}
		fmt.Println("type is", reflect.TypeOf(RecType1))
		fmt.Println("Conn in Server.go ", conn)
		ipaddress, ErrBin := GetIPaddFromConn(conn)
		if !ErrBin {
			log.Fatal("Conn corruption in Server")
		}
		//		fmt.Println("Received Msg", RecType1)
		fmt.Println("Committee Size is", CommitteeSize)
		fmt.Println("IP address is ", ipaddress)

		go conInterface.PassMsgToActor(RecType1, CommitteeSize, ipaddress)

	}

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
