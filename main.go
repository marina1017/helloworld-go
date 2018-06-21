package main

import(
	"io"
	"os"
)

func main() {
	file, err := os.Open("text.txt")
	if err != nil {
		panic(err)
	}
	//「確実に行う後処理」を実行するのに便利な仕組み
	//defer は、現在のスコープが終了したら、その後ろに書かれている行の処理を実行する
	defer file.Close()
	io.Copy(os.Stdout, file)
}