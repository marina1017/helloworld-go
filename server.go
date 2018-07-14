package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	//listen(network名：アドレス)
	//エラーハンドリングしてる
	//複数のゴルーチンは、リスナー上で同時にメソッドを呼び出すことができます。
	//リスナを用いて接続の待ち受け
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running at localhost:8888")

	// 一度で終了しないためにAccept()を何度も繰り返し呼ぶ
	for {
		//Acceptを受け入れて、リスナーへの次の接続を返します。
		//リスナのAcceptメソッドを使用し、クライアントからの接続を待ちます。
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		//goroutine（ゴルーチン）は、Go言語のプログラムで並行に実行されるもの
		// 1リクエスト処理中に他のリクエストのAccept()が行えるように
		// Goroutineを使って非同期にレスポンスを処理する
		// ここではconnを使った読み書きをおこなう
		go func() {
			//RemoteAddrは、リモートネットワークアドレスを返します。
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			// リクエストを読み込む
			//HTTPリクエストのヘッダー、メソッド、パスなどの情報を切り出す
			//NewReaderは、デフォルトのサイズのバッファを持つ新しいReaderを返します。
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			//httputil以下にある便利なデバッグ用の関数
			//リクエスト/レスポンスに含まれるヘッダーとボディのダンプしたバイト列を取得する
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			fmt.Println("dump:\n", string(dump))
			//ProtoMajorとはプロトコルのバージョン番号ぽい　"HTTP/1.0"
			//ioutil 今回は入出力関連のユーティリティ関数が定義されているパッケージ
			//ioutil.NopCloser 関数は、引数の io.Reader に何もしない Close メソッドを付加して
			//io.ReadCloser インターフェースを満たすようにしたオブジェクトを返します
			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body:       ioutil.NopCloser(strings.NewReader("Hello World\n")),
			}
			// レスポンスを書き込む
			response.Write(conn)
			conn.Close()
		}()
	}
}
