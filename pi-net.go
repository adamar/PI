package main

import (
	"bufio"
	"flag"
	"log"
	"os/exec"
	"strings"
)

var pid string

func init() {

	flag.StringVar(&pid, "p", "1", "Pid")
	flag.Parse()

}

func main() {

	cmd := exec.Command("strace", "-f", "-e", "trace=network", "-p", pid)
	stdout, _ := cmd.StderrPipe()
	cmd.Start()
	r := bufio.NewReader(stdout)

	for {

		bufline, err := r.ReadString('\n')
		line := string(bufline)

		val := parseString(line)

			switch val[0]{

			case strings.HasPrefix(lines[3], "recvfrom"):

			case strings.HasPrefix(lines[3], "sendto"):

			case strings.HasPrefix(lines[3], "recvmsg"):

			case strings.HasPrefix(lines[3], "sendmsg"):

			case strings.HasPrefix(lines[3], "getsockopt"):

			case strings.HasPrefix(lines[3], "setsockopt"):

			case strings.HasPrefix(lines[3], "socket"):

			case strings.HasPrefix(lines[3], "connect"):

			case strings.HasPrefix(lines[3], "getsockname"):

			case strings.HasPrefix(lines[3], "bind"):

			}

		}

		if err != nil {
			log.Fatal(err)
		}
	}

}


func parseString(str string) []string {
	w := strings.FieldsFunc(str, func(r rune) bool {
		switch r {
		case '<', '>', ' ', ',', ')', '(', '{', '}':
			return true
		}
		return false
	})
	return w
}

