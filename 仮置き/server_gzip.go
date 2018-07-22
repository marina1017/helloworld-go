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
	"bytes"
	"compress/gzip"
)

// クライアントはgzipを受け入れ可能か？
func isGZipAcceptable(request *http.Request) bool {
	//Indexは、s内でsepが最初に出現する箇所のインデックスを返します。一致しないときは-1を返します。
	//Joinは、パラメータa内の要素を結合して、新たな文字列を作成します。sepで指定したセパレータが結合時に要素間に挿入されます。
	//request.Header["Accept-Encoding"] = gzipが入っている
	//つまりAccept-Encodingのなかみがgzipであれば-1は帰ってこないからtrueになる
	return strings.Index(strings.Join(request.Header["Accept-Encoding"], ","), "gzip") != -1
}

// 1セッションの処理をする
func processSession(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	//defer:特定の処理を関数の一番最後に実行する(第六回で一番最後に書かれていた関数)
	defer conn.Close()

	//keep-alive
	for {
		// タイムアウトを設定
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		// リクエストを読み込む
		//HTTPリクエストのヘッダー、メソッド、パスなどの情報を切り出す
		//NewReaderは、デフォルトのサイズのバッファを持つ新しいReaderを返します。
		request, err := http.ReadRequest(bufio.NewReader(conn))
		print("request:" request)
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

		//　リクエスト表示
		//httputil以下にある便利なデバッグ用の関数
		//リクエスト/レスポンスに含まれるヘッダーとボディのダンプしたバイト列を取得する
		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		// レスポンスを書き込む(ここが変更多い)
		// Header: make(http.Handler)が変更された
		//ここの直接contentを入れるのではなく下でgzipであっしゅくしてからいれる
		response := http.Response{
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: make(http.Header),
		}
		print("make(http.Header):",make(http.Header))
		
		//ここでgzipがつかえるか判定
		if isGZipAcceptable(request) {
			content := "Hello World (gzipped)\n"
			// コンテンツをgzip化して転送
			var buffer bytes.Buffer
			//圧縮した内容はbytes.Bufferに書き出しています。
			
			//定義：func NewWriter(w io.Writer) *Writer { } 引数がio.Writerの構造体 関数内の名前はw　返り値は　Writerのポインタ
			//NewWriterは新しいWriterを返します。
			//返されたライターへの書き込みは圧縮され、wに書き込まれます。
			//完了したら、WriteCloserでCloseを呼び出すのは呼び出し元の責任です。
			//書き込みはバッファリングされ、閉じるまでフラッシュされません。
			//Writer.Headerのフィールドを設定する呼び出し元は、前にそのフィールドを設定する必要があります
			//最初にWrite、Flush、またはCloseを呼び出します。
			//★★疑問：どうしてポインタわたしてるんだろう
			writer := gzip.NewWriter(&buffer)

			//ioI/Oプリミティブへの基本的なインタフェースを提供します。
			//主な役割は、osパッケージ内で定義されているような他のプリミティブを、機能を概念的に表す共通インタフェースへラップすること
			//定義:func WriteString(w Writer, s string) (n int, err os.Error)
			//説明：WriteStringは、wに文字列sの内容を書きこみます。wはバイト配列を受け取ります。
			//★★疑問：あれ？返り値どこいったの？？？？？
			io.WriteString(writer, content)

			//書かれていないデータを基礎となるデータにフラッシュしてライターを閉じます。
			writer.Close()

			//上でつくったレスポンスのbodyに圧縮したcontentをいれる
			//ioutil パッケージ入出力関連のユーティリティ関数が定義されている
			//引数の io.Reader (つまり&buffer)に何もしない Close メソッドを付加して io.ReadCloser インターフェースを満たすようにしたオブジェクトを返します
			response.Body = ioutil.NopCloser(&buffer)

			//Content-Lengthヘッダーに圧縮後のボディサイズを指定します。
			//何故か：６章で紹介されていたHTTP/1.1より前もしくは長さが分からない場合はConnection: closeヘッダーを付与してしまうから
			//本当はKeep-Aliveを付与したい
			response.ContentLength = int64(buffer.Len())

			//クライアント側にgzipで回答してねと伝える
			response.Header.Set("Content-Encoding", "gzip")
		} else {
			content := "Hello World\n"
			response.Body = ioutil.NopCloser(strings.NewReader(content))
			response.ContentLength = int64(len(content))
		}
		response.Write(conn)
	}
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
