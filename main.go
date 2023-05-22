package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("/usr/bin/ffprobe", "-v", "quiet", "-hide_banner", "-loglevel", "warning", "-print_format", "json", "-show_format", "-show-streams", "Lil Durk - All My Life ft. J. Cole.mp4")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error occurred %v", err.Error())
	}

	//j := fmt.Sprintf("%v", out)

	fmt.Printf("\nCmd output:  %v", out.String())
}
