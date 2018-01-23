package gtp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// GTP_Connection GTP连接类管理
type GTPConnection struct {
	cmd     *exec.Cmd
	infile  io.WriteCloser
	outfile io.ReadCloser
}

// NewConnection 创建GTP连接
func NewConnection(cmd string, args ...string) (*GTPConnection, error) {
	conn := GTPConnection{}
	conn.cmd = exec.Command(cmd, args...)
	inf, err := conn.cmd.StdinPipe()
	if err != nil {
		return &conn, err
	}
	outf, err := conn.cmd.StdoutPipe()
	if err != nil {
		return &conn, err
	}
	conn.infile = inf
	conn.outfile = outf
	conn.cmd.Start()
	go func() {
		conn.cmd.Wait()
	}()
	return &conn, nil
}

// Exec 执行GTP命令
func (self GTPConnection) Exec(cmd string) (string, error) {
	self.infile.Write([]byte(fmt.Sprintf("%s \n", cmd)))
	reader := bufio.NewReader(self.outfile)
	result := ""
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		if line == "\n" {
			break
		}
		result += line
	}
	res := strings.Split(result, "")
	l := len(res)
	if res[l-1] == "\n" {
		result = strings.Join(res[:l-1], "")
	}
	if len(result) == 0 {
		return "", errors.New("len =0")
	}
	if res[0] == "?" {
		return "", errors.New(fmt.Sprintf("ERROR: GTP Command failed:%s", strings.Join(res[2:], "")))
	}
	if res[0] == "=" {
		return strings.Join(res[2:], ""), nil
	}
	return "", errors.New(fmt.Sprintf("ERROR: Unrecognized answer: %s", result))
}

// Close 释放GTP资源
func (self GTPConnection) Close() {
	self.infile.Close()
	self.outfile.Close()
}
