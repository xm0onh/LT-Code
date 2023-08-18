package Net

import (
	E "github.com/xm0onh/LT-Code/Encoding"
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

func MsgSender(conn net.Conn, Msg E.VerifyEntity, peer, nodeID, port string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	////	ctx,cancel:=context.WithCancel(context.Background())
	//defer cancel()
	//	fmt.Println("Conn is, ", conn)
	//		fmt.Println("msg is,", Msg)

	//	Encoder := gob.NewEncoder(conn)
	enc := (*MapIdToEncoder)[nodeID]
	err := enc.Encode(&Msg)
	//err := encoder.Encode(&Msg)
	if err != nil {
		fmt.Println("Encoding error is", err.Error())
		conn.Close()
		conn = nil
		time.Sleep(300 * time.Millisecond)
		conn = DialNode(peer, port)
		fmt.Println("Creating new Connection")
		enc := gob.NewEncoder(conn)
		(*IdToConnMap)[nodeID] = conn
		enc = gob.NewEncoder(conn)
		(*MapIdToEncoder)[nodeID] = enc
		MsgSender(conn, Msg, peer, nodeID, port, IdToConnMap, MapIdToEncoder)
		//return enc, true
	}
	//	conn.Close()
	//return enc, false
}

func DialNode(peer, port string) net.Conn {
	conn, err := net.Dial("tcp", peer+":"+port)

	if err != nil {
		fmt.Println(err)
	}
	return conn
}
