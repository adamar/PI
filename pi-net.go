package main

import (
	"bufio"
	"flag"
	"fmt"
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

	fileDescriptors := map[string]string{}
	log.Print(fileDescriptors)

	for {

		bufline, err := r.ReadString('\n')
		line := string(bufline)

		val, empty := parseString(line)

		var prefix string

		if val[0] == "pid" {
			if val[1] != pid {
				prefix = "Child Process[" + val[1] + "] "
			} else {
				prefix = "Parent Process[" + pid + "] "
			}
			val = val[2:]
		}

		//log.Print(val)

		if empty == false {
			switch val[0] {

			case "recvfrom":
				//log.Print(val)
				if val[2] == "unfinished" {

					// recvfrom 9 unfinished ...

					PrintOrange(prefix+"Recvfrom Unfinished", line)
				} else if val[8] == "-1" {

					// recvfrom 9 0x30dfd1340074 4096 0 0 0 = -1 EAGAIN Resource temporarily unavailable

					PrintOrange(prefix+"Failed Recvfrom ", line)
				}

			case "sendto":

			case "recvmsg":

			case "sendmsg":

			case "getsockopt":

			case "setsockopt":

			case "socket":

			case "connect":
				//fd := readFD(val[1], pid, fileDescriptors)
				PrintOrange(prefix+"Connect "+val[7], line)
				//log.Print(fd)

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
		case '<', '>', ' ', ',', ')', '(', '{', '}', '[', ']':
			return true
		}
		return false
	})
	if len(w) < 1 {
		return nil, true
	}
	return w, false
}

func readFD(fd string, pid string, fileDescriptors map[string]string) string {

	if val, ok := fileDescriptors[fd]; ok {
		return val
	}

	out, _ := exec.Command("readlink", "/proc/"+pid+"/fd/"+fd).Output()
	fileDescriptors[fd] = string(out)
	return string(out)

}

func PrintOrange(msg string, msg2 string) {

	fmt.Printf("\x1b[31;1m%s\x1b[0m - %s\n", msg, msg2)

}
