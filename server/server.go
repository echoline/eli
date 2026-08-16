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
	"net"
)

func formatMessage(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9 ]`).ReplaceAllString(s, "")
	return s
}

func learn(rs *rivescript.RiveScript, args []string) string {
	xrs := os.Getenv("HOME") + "/.config/rivescript/learned.rive";

	if len(args) != 2 {
		return "[err: learn args]"
	}

	if len(formatMessage(args[1])) == 0 {
		return "[err: learn: no message found]"
	}

	file, err := os.Open(xrs)
	found := false
	contents := ""
	if err == nil {
		reader := bufio.NewReader(file)
		for {
			m := "+ " + formatMessage(args[1]) + "\n"
			line, err := reader.ReadString('\n')
		        if err != nil && err != io.EOF {
				break
			}

			contents += line
			if line == m {
				found = true
				contents += "- " + args[0] + "\n"
			}

			if err != nil {
				break
			}
		}
		file.Close()
	}

	if found == false {
		contents += "\n+ " + formatMessage(args[1]) + "\n- " + args[0] + "\n"
	}

	data := []byte(contents)
	err = ioutil.WriteFile(xrs, data, 0644)
	if err != nil {
		return "[err: writing to " + xrs + "]"
	}

	rs.LoadFile(xrs)
	rs.SortReplies()

	return ""
}

func main() {
	bot := rivescript.New(nil)

	err := bot.LoadDirectory(os.Getenv("HOME") + "/.config/rivescript/")
	if err != nil {
		fmt.Println("failed to load replies")
		return
	}

	bot.SortReplies()

	bot.SetSubroutine("date", func(rs *rivescript.RiveScript, args []string) string {
		return time.Now().Format("Monday, January 2, 2006")
	})

	bot.SetSubroutine("weekday", func(rs *rivescript.RiveScript, args []string) string {
		return time.Now().Weekday().String()
	})

	bot.SetSubroutine("time", func(rs *rivescript.RiveScript, args []string) string {
		return time.Now().Format("3:04 PM")
	})

	bot.SetSubroutine("hostname", func(rs *rivescript.RiveScript, args []string) string {
		hostname, err := os.Hostname()
		if err != nil {
			return "[err: " + err.Error() + "]"
		}
		return hostname
	})

	os.Remove(os.Getenv("HOME") + "/.config/rivescript.socket")
	listener, err := net.Listen("unix", os.Getenv("HOME") + "/.config/rivescript.socket")
	if err != nil {
		fmt.Printf("unix socket listen failed\n")
		return
	}


	fd, err := listener.Accept()
	if err != nil {
		fmt.Println("unix socket accept failed")
		return
	}

	buf := make([]byte, 8192)
	n, err := fd.Read(buf)
	if err != nil {
		fmt.Println("unix socket read failed")
		return
	}

	text := strings.TrimSpace(string(buf[:n]))
	reply := "";

        for {
		if len(strings.TrimSpace(text)) > 0 {
			reply, _ = bot.Reply(os.Getenv("USER"), text)
		} else {
			reply = ""
		}

		n, err = fd.Write([]byte(reply))
		if err != nil {
			fmt.Println("unix socket write failed")
			return
		}

		fd.Close()

		fd, err := listener.Accept()
		if err != nil {
			fmt.Println("unix socket accept failed")
			return
		}

		n, err := fd.Read(buf)
		if err != nil {
			fmt.Println("unix socket read failed")
			return
                }

		text = strings.TrimSpace(string(buf[:n]))

		if (len(text) > 0 && len(reply) > 0) {
			args := []string{text, reply}
			learn(bot, args)
		}
	}
}

