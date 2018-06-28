package main

import(
	"io"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader("Example of io.SectionReader\n")
	//Section の部分だけを切り出した Reader をまず作成し
	sectionReader := io.NewSectionReader(reader, 14, 7)
	//それをすべて os.Stdout に書き出しています
	io.Copy(os.Stdout, sectionReader)
}