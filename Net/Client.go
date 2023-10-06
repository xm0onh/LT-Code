package Net

import (
	E "github.com/xm0onh/LT-Code/Encoding"
	kzg "github.com/xm0onh/LT-Code/KZG"

	//	"github.com/xm0onh/LT-Code/Timer"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func MsgSender(conn net.Conn, Msg E.VerifyEntity, peer, nodeID, port string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	enc := (*MapIdToEncoder)[nodeID]

	dataType := "VerifyEntity"
	err := enc.Encode(&dataType)
	if err != nil {
		handleEncodingError(err, conn, peer, port, Msg, nodeID, IdToConnMap, MapIdToEncoder)
		return
	}

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
	dataType := "KZGZSender"
	err := enc.Encode(&dataType)
	if err != nil {
		handleEncodingErrorZ(err, conn, peer, port, Z, nodeID, IdToConnMap, MapIdToEncoder)
		return
	}

	err = enc.Encode(&Z)
	if err != nil {
		handleEncodingErrorZ(err, conn, peer, port, Z, nodeID, IdToConnMap, MapIdToEncoder)
	}
}

func handleEncodingErrorZ(err error, conn net.Conn, peer, port string, Z kzg.KZGZSender, nodeID string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	fmt.Println("Encoding error:", err)
	conn.Close()
	time.Sleep(300 * time.Millisecond)
	newConn := DialNode(peer, port)
	fmt.Println("Creating new Connection")
	newEnc := gob.NewEncoder(newConn)
	(*IdToConnMap)[nodeID] = newConn
	(*MapIdToEncoder)[nodeID] = newEnc
	KZGZSender(newConn, Z, peer, nodeID, port, IdToConnMap, MapIdToEncoder)
}

func KZGZVerifier(conn net.Conn, Z kzg.KZGVerifier, peer, nodeID, port string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	enc := (*MapIdToEncoder)[nodeID]
	dataType := "KZGZVerifier"
	err := enc.Encode(&dataType)
	if err != nil {
		handleEncodingErrorKZGVerify(err, conn, peer, port, Z, nodeID, IdToConnMap, MapIdToEncoder)
		return
	}

	err = enc.Encode(&Z)
	if err != nil {
		handleEncodingErrorKZGVerify(err, conn, peer, port, Z, nodeID, IdToConnMap, MapIdToEncoder)
	}
}

func handleEncodingErrorKZGVerify(err error, conn net.Conn, peer, port string, Z kzg.KZGVerifier, nodeID string, IdToConnMap *map[string]net.Conn, MapIdToEncoder *map[string]*gob.Encoder) {
	fmt.Println("Encoding error:", err)
	conn.Close()
	time.Sleep(300 * time.Millisecond)
	newConn := DialNode(peer, port)
	fmt.Println("Creating new Connection")
	newEnc := gob.NewEncoder(newConn)
	(*IdToConnMap)[nodeID] = newConn
	(*MapIdToEncoder)[nodeID] = newEnc
	KZGZVerifier(newConn, Z, peer, nodeID, port, IdToConnMap, MapIdToEncoder)
}

func DialNode(peer, port string) net.Conn {
	fmt.Println("tcp", peer+":"+port)
	conn, err := net.Dial("tcp", peer+":"+port)

	if err != nil {
		fmt.Println(err)
	}
	return conn
}
