package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {
	var err error
		//まだコネクションを張っていない/エラーでリトライ時はDialから行う
		if conn == nil {
			conn, err = net.Dial("tcp", "localhost:8888")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Access: %d\n", current)
		}

		//POSTで文字列を送るリクエストを作成
		request, err := http.NewRequest("POST", "http://localhost:8888", nil)
		if err != nil {
			panic(err)
		}
		//サーバー側にgzip対応してるか聞く
		request.Header.Set("Accept-Encoding", "gzip")

		//リクエストを送る
		err = request.Write(conn)
		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(conn)

		//サーバーから読み込む　タイムアウトはここでエラーになるのでリトライ
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			fmt.Println("retry")
			conn = nil
			continue
		}

		//結果を表示
		//定義： func DumpResponse(resp *http.Response, body bool) ([]byte, error) {
		//第一引数が*http.Response型 第二引数はbodyを読み込むかをbool値できめることができる　返り値は([]byte, error)
		//６章とちがってbodyはgzipで圧縮されているのdumpResponse関数で読めないため無視

		dump, err := httputil.DumpResponse(response, false)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		//サーバーから帰ってくるレスポンスにTransferEncodingが入っててそれがchunkedかチェック
		if len(response.TransferEncoding) < 1 || response.TransferEncoding[0] != "chunked" {
			panic("wrong transfer encoding")
		}

		//forループでチャンクごとに読み込んでいます。
		for {
			// 改行をさがしてサイズを取得
			sizeStr, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			// 16進数のサイズをパース。サイズがゼロならクローズ
			//定義:func ParseInt(s string, base int, bitSize int) (i int64, err error)
			//文字列を任意の基数(2進数〜36進数)・任意のビット長(8〜64bit)のIntにパースする。
			//16進数で64ビットのIntにパースする
			size, err := strconv.ParseInt(string(sizeStr[:len(sizeStr)-2]), 16, 64)
			if size == 0 {
				break
			}
			if err != nil {
				panic(err)
			}
			// サイズ数分バッファを確保して読み込み
			//go channelはgoroutineの間で値を受け渡しするための配列のようなものです。
    		line := make([]byte, int(size))
			//このreaderとは上で定義されているreader := bufio.NewReader(conn)のこと
    		reader.Read(line)
			//Discardは、次のnバイトをスキップし、破棄されたバイト数を返します。
			//★★疑問：この辺の返り値はどこに行ったの・・・？　この2って改行処理の\2バイト分ってことかな？
    		reader.Discard(2)
    		fmt.Printf("  %d bytes: %s\n", size, string(line))
  		}
}
