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

	cmd := exec.Command("sudo", "strace", "-f", "-e", "trace=network", "-p", pid)
	stdout, _ := cmd.StderrPipe()
	cmd.Start()
	r := bufio.NewReader(stdout)

	fileDescriptors := map[string]string{}

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

		if empty == false {
			switch val[0] {
			case "recvfrom":
				fd := readFD(val[1], pid, fileDescriptors)
				parseRecvfrom(prefix, val, fd)
			case "sendto":
				fd := readFD(val[1], pid, fileDescriptors)
				parseSendto(prefix, val, fd)
			case "recvmsg":
				fd := readFD(val[1], pid, fileDescriptors)
				parseRecvmsg(prefix, val, fd)
			case "sendmsg":
				fd := readFD(val[1], pid, fileDescriptors)
				parseSendmsg(prefix, val, fd)

			//case "getsockopt":

			//case "setsockopt":

			//case "socket":

			case "connect":
				parseConnect(prefix, val)

			//case "getsockname":

			//case "bind":

			default:
				log.Print(val)
			}

		}
		if err != nil {
			log.Fatal(err)
		}
	}

}

func parseRecvmsg(prefix string, val []string, fd string) {

	PrintOrange(prefix + "Recvmsg from " + fd)

}

func parseSendmsg(prefix string, val []string, fd string) {

	PrintOrange(prefix + "Sendmsg to " + fd)

}

func parseSendto(prefix string, val []string, fd string) {

	PrintOrange(prefix + "Send to " + fd)

}

func parseConnect(prefix string, val []string) {

	PrintOrange(prefix + "Connect " + val[7])

}

func parseRecvfrom(prefix string, val []string, fd string) {

	// recvfrom 9 unfinished ...
	if val[2] == "unfinished" {

		PrintOrange(prefix + "Recvfrom Unfinished " + fd)

		// recvfrom 9 0x30dfd1340074 4096 0 0 0 = -1 EAGAIN Resource temporarily unavailable
	} else if val[8] == "-1" {

		PrintOrange(prefix + "Failed Recvfrom " + fd)

	} else {

		log.Print(val)

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

func PrintOrange(msg string) {

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", msg)

}
