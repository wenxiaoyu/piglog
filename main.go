// main.go
package main

import (
	"fmt"
	piglog "piglog/libs"
)

func main() {
	fmt.Println("piglog sever start... ")

	if piglog.LogConfig.Remote.Ip != "" && piglog.LogConfig.Remote.Port > 0 {
		fmt.Print("start http service : ")
		fmt.Println(piglog.LogConfig.Remote)
		piglog.Log.StartHttpServer(piglog.LogConfig.Remote)
	}
}
