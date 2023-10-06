package Net

import (
	"reflect"

	E "github.com/xm0onh/LT-Code/Encoding"
	kzg "github.com/xm0onh/LT-Code/KZG"

	//	"github.com/xm0onh/LT-Code/Timer"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

//func BroadCastMsg(Msg E.VerifyEntity, peers []string, Msgport string, connSlice []net.Conn){
//	//fmt.Println("Msg is,", Msg)
//	//dialSlice:=dial(peers,Msgport)
//	for indx,conn:=range connSlice{
//		go MsgSender(conn,Msg, peers[indx], Msgport)
//	}

//}

// func MsgSender(conn net.Conn, Msg E.VerifyEntity, peer, nodeID, port string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
// 	////	ctx,cancel:=context.WithCancel(context.Background())
// 	//defer cancel()
// 	//	fmt.Println("Conn is, ", conn)
// 	//		fmt.Println("msg is,", Msg)

// 	//	Encoder := gob.NewEncoder(conn)
// 	enc := (*MapIdToEncoder)[nodeID]
// 	// fmt.Println("Encoder is", enc)
// 	// fmt.Println("Encoder type is", reflect.TypeOf(enc))
// 	// fmt.Println("Msg type is", reflect.TypeOf(Msg))
// 	err := enc.Encode(&Msg)
// 	// err := encoder.Encode(&Msg)
// 	if err != nil {
// 		fmt.Println("Encoding error is", err.Error())
// 		conn.Close()
// 		conn = nil
// 		time.Sleep(300 * time.Millisecond)
// 		conn = DialNode(peer, port)
// 		fmt.Println("Creating new Connection")
// 		enc := gob.NewEncoder(conn)
// 		(*IdToConnMap)[nodeID] = conn
// 		(*MapIdToEncoder)[nodeID] = enc
// 		MsgSender(conn, Msg, peer, nodeID, port, IdToConnMap, MapIdToEncoder)
// 		//return enc, true
// 	}
// 	//	conn.Close()
// 	//return enc, false
// }

func MsgSender(conn net.Conn, Msg E.VerifyEntity, peer, nodeID, port string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	enc := (*MapIdToEncoder)[nodeID]

	// First send the dataType
	dataType := "VerifyEntity" // This is the data type identifier for E.VerifyEntity
	err := enc.Encode(&dataType)
	if err != nil {
		handleEncodingError(err, conn, peer, port, Msg, nodeID, IdToConnMap, MapIdToEncoder)
		return
	}

	// Then send the actual message
	err = enc.Encode(&Msg)
	if err != nil {
		handleEncodingError(err, conn, peer, port, Msg, nodeID, IdToConnMap, MapIdToEncoder)
	}
}

func handleEncodingError(err error, conn net.Conn, peer, port string, Msg E.VerifyEntity, nodeID string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	fmt.Println("Encoding error:", err)
	conn.Close()
	time.Sleep(300 * time.Millisecond)
	newConn := DialNode(peer, port)
	fmt.Println("Creating new Connection")
	newEnc := gob.NewEncoder(newConn)
	(*IdToConnMap)[nodeID] = newConn
	(*MapIdToEncoder)[nodeID] = newEnc
	MsgSender(newConn, Msg, peer, nodeID, port, IdToConnMap, MapIdToEncoder)
}

func KZGZSender(conn net.Conn, Z kzg.KZGZSender, peer, nodeID, port string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	enc := (*MapIdToEncoder)[nodeID]
	// fmt.Println("Encoder is", enc)
	// fmt.Println("Encoder type is", reflect.TypeOf(enc))
	fmt.Println("Msg type is", reflect.TypeOf(Z))
	err := enc.Encode(&Z)
	fmt.Println(err)
	// err := encoder.Encode(&Msg)
	if err != nil {
		fmt.Println("Encoding error is", err.Error())
		conn.Close()
		conn = nil
		time.Sleep(300 * time.Millisecond)
		conn = DialNode(peer, port)
		fmt.Println("Creating new Connection")
		enc := gob.NewEncoder(conn)
		(*IdToConnMap)[nodeID] = conn
		(*MapIdToEncoder)[nodeID] = enc
		KZGZSender(conn, Z, peer, nodeID, port, IdToConnMap, MapIdToEncoder)
		//return enc, true
	}
	//	conn.Close()
	//return enc, false
}

func DialNode(peer, port string) net.Conn {
	fmt.Println("tcp", peer+":"+port)
	conn, err := net.Dial("tcp", peer+":"+port)

	if err != nil {
		fmt.Println(err)
	}
	return conn
}
