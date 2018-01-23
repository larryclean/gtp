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
	conn, err := gtp.NewConnection("./gnugo_mac", "--mode", "gtp")
	checkError(err)
	client:=gtp.NewGtpClient(conn)
	move,err:=client.GenMove("B")
	fmt.Println(move,err)
	move,err=client.GenMove("W")
	fmt.Println(move,err)
	board,err:=client.ShowBoard()
	fmt.Println(board,err)
}
