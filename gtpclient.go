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
	if strings.Contains(ver, "ERROR") {
		util.protocol_version = ver
	} else {
		util.protocol_version = "1"
	}
	return &util
}

func (self GTPClient) KnowCommand(cmd string) (string, error) {
	return self.conn.Exec(cmd)
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
func (self GTPClient) ShowBoard() (string, error) {
	return self.conn.Exec("showboard")
}
func (self GTPClient) ClearBoard() (string, error) {
	return self.conn.Exec("clear_board")
}

func (self GTPClient) FinalScore() (string, error) {
	return self.conn.Exec("final_score")
}
func (self GTPClient) Quit() (string, error) {
	return self.conn.Exec("Quit")
}
