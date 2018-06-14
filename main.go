package main

import(
    "os"
)

func main() {
    // os.Create新規ファイルの作成
    file, err:=os.Create("text.txt")
    if err!=nil{
        panic(err)
    }
    //Write が受け取るのは文字列ではなくてバイト列だからバイト列に変換
    //定義　func (f *File) Write(b []byte) (n int, err error) {
    file.Write([]byte("os.File example\n"))
    file.Close()
}
