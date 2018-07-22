package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
	//このふたつはいった
)

// 順番に従ってconnに書き出しをする(goroutineで実行される)（右から出てくるぶぶん）
func writeToConn(sessionResponses chan chan *http.Response, conn net.Conn) {
	//defer最後に実行
	defer conn.Close()
	// 順番に取り出す
	for sessionResponse := range sessionResponses {

		// 選択された仕事が終わるまで待つ
		//データの入出力には<-演算子を使います
		// データ投入
		//buffered <- "データ"
		// データ取り出し
		//variable <- buffered
		response := <-sessionResponse
		response.Write(conn)
		close(sessionResponse)
	}
}

// リクエストごとに非同期処理でレスポンスを返す処理　セッション内のリクエストを処理する（左に入れるぶぶん）
func handleRequest(request *http.Request, resultReceiver chan *http.Response) {
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
	content := "Hello World\n"
	// レスポンスを書き込む
	// セッションを維持するためにKeep-Aliveでないといけない
	response := &http.Response{
		StatusCode:    200,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: int64(len(content)),
		Body:          ioutil.NopCloser(strings.NewReader(content)),
	}
	// 処理が終わったらチャネルに書き込み、
	// ブロックされていたwriteToConnの処理を再始動する
	// データ投入
	//buffered <- "データ"
	resultReceiver <- response
}

// セッション1つを処理
func processSession(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	// セッション内のリクエストを順に処理するためのチャネル
	// レスポンスの順番を制御するためには、Go言語のデータ構造のチャネルを使っています
	//チャネルはFIFO(First In, First Outを表す頭字語である。 先入れ先出し)のキューで、バッファなしとバッファありの2種類があります。
	//例
	////////////////////////////////////
	// バッファなし
	//unbuffered := make(chan string)
	// バッファあり
	//buffered := make(chan string, 10)
	////////////////////////////////////

	//バッファありの場合は、指定した個数までは自由に投入できますが、
	// 指定した個数のデータが入っているときにさらに追加でデータを投入しようとすると、投入しようとしたスレッド（ゴルーチン）がブロックされます。
	//50バッファあり
	sessionResponses := make(chan chan *http.Response, 50)

	//defer　最後に処理
	defer close(sessionResponses)

	// レスポンスを直列化してソケットに書き出す専用のゴルーチン平行処理(上に飛ぶ)
	go writeToConn(sessionResponses, conn)

	////左からはいってくるのがここから
	reader := bufio.NewReader(conn)
	for {
		// レスポンスを受け取ってセッションのキューに
		// 入れる
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		// リクエストを読み込む
		request, err := http.ReadRequest(reader)
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				fmt.Println("Timeout")
				break
			} else if err == io.EOF {
				break
			}
			panic(err)
		}
		// バッファなし
		//unbuffered := make(chan string)
		sessionResponse := make(chan *http.Response)

		// 上の５０バッファがあるsessionResponsesにバッファなしのsessionResponseデータ投入
		//buffered <- "データ"
		sessionResponses <- sessionResponse

		// 非同期でレスポンスを実行(左側からはいってくる)
		go handleRequest(request, sessionResponse)
	}
	////左からはいってくるのがここまで
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running at localhost:8888")
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		//スレッドよりも小さい処理単位である「ゴルーチン(goroutine)」が並行して動作するように実装されている。
		//go文は、このゴルーチンを新たに生成して、並行して処理される新処理の流れをランタイムに追加するための機能。
		go processSession(conn)
	}
}
