package assignment03IBC

import (
	"bufio"
	"encoding/gob"
	"fmt"
	a2 "github.com/shaheerhas/assignment02IBC"
	"log"
	"net"
)

type node struct {
	connection net.Conn
	address string
}
var Quorum int
var connections[] node
var ChainHead *a2.Block

func StartListening( myListeningAddress string,  listener string){
	if listener == "Satoshi" {
		ChainHead = a2.InsertBlock("","",listener,0,ChainHead)

		ln, err := net.Listen("tcp", myListeningAddress)
		if err != nil {
			log.Fatal(err)
		}
		for {
			var Node node
			conn, err := ln.Accept()
			Node.connection = conn

			connections = append(connections, Node)

			if err != nil {
				log.Println(err)
				continue
			}
			go handleConnection(conn, listener)
		}
	}else{

		ln, err := net.Listen("tcp", myListeningAddress)
		if err != nil {
			log.Fatal(err)
		}

			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}

		fmt.Println(ReadString(conn))

	}

}
func handleConnection(conn net.Conn,listener string){
	Quorum--
	i:=0
	for i < len(connections){
		if connections[i].connection == conn {
			connections[i].address = ReadString(conn)
			break
		}
	}
	ChainHead = a2.InsertBlock("","",listener,0,ChainHead)
}
func WriteString(conn net.Conn, myListeningAddress string){
	conn.Write([]byte(myListeningAddress+"\n"))
}
func ReadString(conn net.Conn) string {

	clientReader := bufio.NewReader(conn)
	str, _, _ := clientReader.ReadLine()
	return string(str)
}
func ReceiveChain(conn net.Conn) a2.Block{

	var recvdBlock a2.Block
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&recvdBlock)
	if err != nil {
		//handle error
		fmt.Println("decode error")
	}
	return recvdBlock
}
func WaitForQuorum(){
	for Quorum > 0{
	}
}
func SendChainandConnInfo(){
	i:=0
	for i < len(connections) {
		gobEncoder := gob.NewEncoder(connections[i].connection)
		err := gobEncoder.Encode(ChainHead)
		if err != nil {
			log.Println(err)
		}

		if i==0{
			WriteString(connections[i].connection,connections[len(connections)-1].address)
		}else {
			WriteString(connections[i].connection,connections[i-1].address)
		}
		i++
	}
}
