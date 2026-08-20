package main

import (
	"fmt"
	"bufio"
	"os"
	"net"
	"time"
	"strings"
)

func wrapText(text string, limit int) string {
	prompt := "< ";

	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(prompt + words[0])
	spaceLeft := limit - len(words[0]) - len(prompt) 

	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			result.WriteString("\n" + prompt + word)
			spaceLeft = limit - len(word) - len(prompt)
		} else {
			result.WriteString(" " + word)
			spaceLeft -= len(word) + 1
		}
	}

	return result.String()
}

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

	fmt.Println(wrapText(string(buf[:n]), 80))
}

