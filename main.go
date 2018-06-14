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
    //Write が受け取るのは文字列ではなくてバイト列
    file.Write([]byte("os.File example\n"))
    file.Close()
}
