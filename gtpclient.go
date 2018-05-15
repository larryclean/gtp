package gtp

import (
	"fmt"
	"strings"
)

type GTPClient struct {
	conn             *GTPConnection
	protocol_version string
}

func NewGtpClient(conn *GTPConnection) *GTPClient {
	util := GTPClient{}
	util.conn = conn
	ver, _ := conn.Exec("protocol_version")
	fmt.Println(111)
	if strings.Contains(ver, "ERROR") {
		util.protocol_version = ver
	} else {
		util.protocol_version = "2"
	}
	return &util
}

func (self GTPClient) KnowCommand(cmd string) (string, error) {
	return self.conn.Exec("known_command "+cmd)
}

func (self GTPClient) GenMove(color string) (string, error) {
	color = strings.ToUpper(color)
	command := "black"
	if color == "B" {
		command = "black"
	} else if color == "W" {
		command = "white"
	}
	if self.protocol_version == "1" {
		command = "genmove_" + command
	} else {
		command = "genmove " + command
	}
	return self.conn.Exec(command)
}
func (self GTPClient) Move(color, coor string) (string, error) {
	color = strings.ToUpper(color)
	command := "black"
	if color == "B" {
		command = "black"
	} else if color == "W" {
		command = "white"
	}
	return self.conn.Exec(fmt.Sprintf("play %s %s", command, coor))
}

func (self GTPClient) LoadSgf(file string, move int) (string, error) {
	command := fmt.Sprintf("loadsgf %s", file)
	return self.conn.Exec(command)
}
func (self GTPClient) FinalStatusList(cmd string) (string, error) {
	command := fmt.Sprintf("final_status_list %s", cmd)
	return self.conn.Exec(command)
}
func (self GTPClient) SetLevel(seed int) (string, error) {
	command := fmt.Sprintf("level %d", seed)
	return self.conn.Exec(command)
}
func (self GTPClient) SetRandomSeed(seed int) (string, error) {
	command := fmt.Sprintf("set_random_seed %d", seed)
	return self.conn.Exec(command)
}
func (self GTPClient) BoardSize(size int) (string, error) {
	command := fmt.Sprintf("boardsize %d", size)
	return self.conn.Exec(command)
}
func (self GTPClient) ShowBoard() (string, error) {
	return self.conn.Exec("showboard")
}
func (self GTPClient) ClearBoard() (string, error) {
	return self.conn.Exec("clear_board")
}
func (self GTPClient) PrintSgf() (string, error) {
	return self.conn.Exec("printsgf")
}
func (self GTPClient) TimeSetting(baseTime, byoTime, byoStones int) (string, error) {
	return self.conn.Exec(fmt.Sprintf("time_settings %d %d %d", baseTime, byoTime, byoStones))
}
func (self GTPClient) FinalScore() (string, error) {
	return self.conn.Exec("final_score")
}
func (self GTPClient) Quit() (string, error) {
	return self.conn.Exec("Quit")
}

func (self GTPClient) SendCmd(cmd string) (string, error) {
	return self.conn.Exec(cmd)
}
