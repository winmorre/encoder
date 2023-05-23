package main

import (
	"fmt"
	"os/exec"
)

func main() {
	out, err := exec.Command("ffprobe", "-version").CombinedOutput()
	//var out strings.Builder
	//cmd.Stdout = &out
	//err := cmd.Run()

	if err != nil {
		fmt.Printf("Error occurred %v", err.Error())
	}

	//j := fmt.Sprintf("%v", out)

	fmt.Printf("\nCmd output:  %v", string(out))
}
