package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Server is running at localhost:8888")
	//UDPが使うのがListenPacket
	conn, err := net.ListenPacket("udp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	//データサイズが決まらないデータに対しては、
	//フレームサイズ分のバッファや、期待されるデータの最大サイズ分のバッファを作り、そこにデータをまとめて読み込むことになります。
	buffer := make([]byte, 1500)
	for {
		//ReadFrom()メソッドを使うと、通信内容を読み込むと同時に、接続してきた相手のアドレス情報が受け取れます。
		length, remoteAddress, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received from %v: %v\n", remoteAddress, string(buffer[:length]))
		_, err = conn.WriteTo([]byte("Hello from Server"), remoteAddress)
		if err != nil {
			panic(err)
		}
	}
}
