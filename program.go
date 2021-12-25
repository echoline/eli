package main

import (
	"github.com/aichaos/rivescript-go"
	"time"
	"fmt"
	"strings"
	"bufio"
	"os"
	"io/ioutil"
	"io"
	"regexp"
)

func formatMessage(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9 ]`).ReplaceAllString(s, "")
	return s
}

func main() {
	bot := rivescript.New(nil)
	reader := bufio.NewReader(os.Stdin)

	err := bot.LoadDirectory("./replies/")
	if err != nil {
		fmt.Printf("failed to load replies\n")
		return
	}

	bot.SortReplies()

	fmt.Println("loaded replies")

	bot.SetSubroutine("time", func(rs *rivescript.RiveScript, args []string) string {
		return time.Now().Format(time.RFC1123)
	})

	bot.SetSubroutine("today", func(rs *rivescript.RiveScript, args []string) string {
		return time.Now().Weekday().String()
	})

	bot.SetSubroutine("learn", func(rs *rivescript.RiveScript, args []string) string {
		xrs := args[0]
		s := strings.Split(strings.Join(args[1:], " "), "::::")
		if len(s) >= 2 {
			if len(formatMessage(s[1])) == 0 {
				return "[err: no message found]"
			}
			file, err := os.Open(xrs)
			found := false
			contents := ""
			if err == nil {
				reader := bufio.NewReader(file)
				for {
					m := "+ " + formatMessage(s[1]) + "\n"
					line, err := reader.ReadString('\n')
				        if err != nil && err != io.EOF {
						break
					}

					contents += line
					if line == m {
						found = true
						contents += "- " + s[0] + "\n"
					}

					if err != nil {
						break
					}
				}
				file.Close()
			}
			if found == false {
				contents += "\n+ " + formatMessage(s[1]) + "\n- " + s[0] + "\n"
			}
			data := []byte(contents)
			err = ioutil.WriteFile(xrs, data, 0644)
			if err != nil {
				return "error writing to " + xrs
			}
			bot.LoadFile(xrs)
			bot.SortReplies()
			if len(s) == 3 {
				return ""
			}
			return "Okay, I'll try to remember to respond, \"" + s[0] + "\" when you say, \"" + s[1] + "\""
		}
		return ""
	})

	for {
		text, _ := reader.ReadString('\n')

		if len(strings.TrimSpace(text)) > 0 {
			reply, _ := bot.Reply("client username goes here", text)

			fmt.Println(reply)
		}
	}
}

