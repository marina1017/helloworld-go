package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listen tick server at 224.0.0.1:9999")
	//net.ResolveUDPAddr()関数でパースしてから
	address, err := net.ResolveUDPAddr("udp", "224.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	//TCPのときに使ったnet.Listen()ではなく UDPによるマルチキャスト専用のnet.ListenMulticastUDP()という関数
	listener, err := net.ListenMulticastUDP("udp", nil, address)
	defer listener.Close()
	buffer := make([]byte, 1500)
	for {
		length, remoteAddress, err := listener.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Server %v\n", remoteAddress)
		fmt.Printf("Now    %s\n", string(buffer[:length]))
	}
}
