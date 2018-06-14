package main
import(
	"bytes"
	"fmt"
)
func main() {
	var buffer bytes.Buffer
	buffer.WriteString("bytes.Buffer example\n")
	fmt.Println(buffer.String())
}