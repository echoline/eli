package main

import (
	"fmt"
	"bufio"
	"os"
	"net"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	socket, err := net.DialTimeout("unix", os.Getenv("HOME") + "/.config/rivescript.socket", time.Second)
	if err != nil {
		fmt.Println("error connecting to unix socket: " + err.Error())
	}

	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading from stdin: " + err.Error())
	}

	n, err := socket.Write([]byte(text))
	if err != nil {
		fmt.Println("error writing to unix socket: " + err.Error())
	}

	buf := make([]byte, 8192)
	n, err = socket.Read(buf)
	if err != nil {
		fmt.Println("error reading from unix socket: " + err.Error())
	}

	fmt.Println(string(buf[:n]))
}

