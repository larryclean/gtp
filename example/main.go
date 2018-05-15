package main

import (
	"github.com/larry-dev/gtp"
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
func main() {
	conn, err := gtp.NewConnection("./leelaz.mac.0.15", "-g", "-w","/Users/larryliu/Documents/01-Project/gopath/src/github.com/larry-dev/gtp/example/network")
	checkError(err)
	client:=gtp.NewGtpClient(conn)
	value:=client.KnowCommand("kgs-time_settings")
	if value{
		fmt.Println("支持kgs-time_settings")
	}
	move,err:=client.GenMove("B")
	fmt.Println(move,err)
	move,err=client.GenMove("W")
	fmt.Println(move,err)
	board,err:=client.ShowBoard()
	fmt.Println(board,err)
}
