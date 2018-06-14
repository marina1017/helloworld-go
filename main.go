package main

import(
    "os"
)

func main() {
    //前回の fmt.Println では、最終的に os.Stdout の Write メソッドを呼び出していました。
    os.Stdout.Write([]byte("os.Stdout example\n"))
}
