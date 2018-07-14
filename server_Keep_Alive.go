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
			// Accept後のソケットで何度も応答を返すためにループ
			for {
				// タイムアウトを設定
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))

				// リクエストを読み込む
				//HTTPリクエストのヘッダー、メソッド、パスなどの情報を切り出す
				//NewReaderは、デフォルトのサイズのバッファを持つ新しいReaderを返します。
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					// タイムアウトもしくはソケットクローズ時は終了
					// それ以外はエラーにする
					neterr, ok := err.(net.Error) // ネットワークのエラーに変換している
					if ok && neterr.Timeout() {
						fmt.Println("Timeout")
						break
						//EOFとは「ここでファイルは終わりですよ～」を表す目印
					} else if err == io.EOF {
						break
					}
					panic(err)
				}

				//　リクエスト表示
				//httputil以下にある便利なデバッグ用の関数
				//リクエスト/レスポンスに含まれるヘッダーとボディのダンプしたバイト列を取得する
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println("dump:\n", string(dump))
				content := "Hello World\n"

				// レスポンスを書き込む
				// HTTP/1.1かつ、ContentLengthの設定が必要
				//Go言語のResponse.Write()は、HTTP/1.1より前もしくは長さが分からない場合は
				//Connection: closeヘッダーを付与してしまうから設定が必要
				response := http.Response{
					StatusCode:    200,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)),
					Body:          ioutil.NopCloser(strings.NewReader(content)),
				}
				response.Write(conn)

			}
			conn.Close()
		}()
	}
}
