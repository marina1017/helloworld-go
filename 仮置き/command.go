package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) == 1 {
		return
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	state := cmd.ProcessState
	// 終了コードと状態を文字列で返す
	fmt.Printf("%s\n", state.String())
	// 子プロセスのプロセスID
	fmt.Printf("  Pid: %d\n", state.Pid())
	// 終了しているかどうか
	fmt.Printf("  Exited: %v\n", state.Exited())
	// 正常終了か？
	fmt.Printf("  Success: %v\n", state.Success())
	fmt.Printf("%s\n", state.String())
	fmt.Printf("  Pid: %d\n", state.Pid())
	fmt.Printf("  System: %v\n", state.SystemTime())
	fmt.Printf("  User: %v\n", state.UserTime())

	env := cmd.SysProcAttr
	fmt.Printf("Success: %v\n", env.ExitStatus())
}
