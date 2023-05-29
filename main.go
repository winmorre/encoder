package main

import (
	"encoder/routes"
	"fmt"
)

func main() {
	//out, err := exec.Command("ffprobe", "-version").CombinedOutput()
	//var out strings.Builder
	//cmd.Stdout = &out
	//err := cmd.Run()

	//if err != nil {
	//	fmt.Printf("Error occurred %v", err.Error())
	//}

	//j := fmt.Sprintf("%v", out)

	//fmt.Printf("\nCmd output:  %v", string(out))
	runErr := routes.RegisterRoutes().Run()
	if runErr != nil {
		fmt.Printf("Error Occurred running server, %v", runErr.Error())
		return
	}
}
