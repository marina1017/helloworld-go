package main

import(
	"os"
	"fmt"
	"time"
)

func main() {
	fmt.Fprintf(os.Stdout, "Write with os.Stdout at %v", time.Now())
}