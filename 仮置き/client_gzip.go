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
	sendMessage := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}
	current := 0
	var conn net.Conn = nil

	//リトライ用にループで全体を囲う
	for {
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
		//増えた第三引数はbodyに入る
		//NewRequest(method, url string, body io.Reader)返り値は(*Request, error)のふたつ
		request, err := http.NewRequest("POST", "http://localhost:8888", strings.NewReader(sendMessage[current]))
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
		//deferは最後に実行する関数
		//疑問最後にかけばいいのになんでわざわざここにかくの？
		defer response.Body.Close()

		//サーバーから送られてきたレスポンスのbodyの中身を確認 gzipで圧縮されてるかみる
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body)
			if err != nil {
				panic(err)
			}
			//このCopy関数でgzipで圧縮された内容Hello World (gzipped)を出力している
			//疑問:でもなんでPrint関数でもないのにCopy関数で出力できるのかはよくわからない
			io.Copy(os.Stdout, reader)
			reader.Close()
		} else {
			io.Copy(os.Stdout, response.Body)
		}

		//何個送っておいたか保持しておく
		current++
		//sendMessageにはいっている"ASCII","PROGRAMMING","PLUS",の数とcurrentの数が一致したら終了
		if current == len(sendMessage) {
			break
		}
	}
}
