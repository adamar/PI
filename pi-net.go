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

		val, empty := parseString(line)

		if empty == false {
			switch val[0] {

			case "recvfrom":

			case "sendto":

			case "recvmsg":

			case "sendmsg":

			case "getsockopt":

			case "setsockopt":

			case "socket":

			case "connect":

			case "getsockname":

			case "bind":

			}

		}
		if err != nil {
			log.Fatal(err)
		}
	}

}

func parseString(str string) ([]string, bool) {
	w := strings.FieldsFunc(str, func(r rune) bool {
		switch r {
		case '<', '>', ' ', ',', ')', '(', '{', '}':
			return true
		}
		return false
	})
	if len(w) < 1 {
		return nil, true
	}
	return w, false
}
