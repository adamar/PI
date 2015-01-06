package main

import (
	"bufio"
	"log"
	"os/exec"
	"flag"
	"strings"
)

var pid string

func init() {

	flag.StringVar(&pid, "p", "1", "Pid")
	flag.Parse()

}

func main() {

	cmd := exec.Command("strace", "-f", "-e", "trace=file", "-p", pid)
	stdout, _ := cmd.StderrPipe()
	cmd.Start()
	r := bufio.NewReader(stdout)

	for {

		bufline, err := r.ReadString('\n')
		line := string(bufline)

		lines := strings.Split(line, " ")
		if len(lines) > 3 {
			filestring := strings.Split(lines[3], `"`)

			switch {

			case strings.HasPrefix(lines[3], "open"):
                print("open " + filestring[1])

			case strings.HasPrefix(lines[3], "stat"):
				print("stat " + filestring[1])

			case strings.HasPrefix(lines[3], "readlink"):
				print("check symlink " + filestring[1])

			default:
				print("Undefined :" + line)

			}
		}

		if err != nil {
			log.Fatal(err)
		}
	}

}

func print(txt string) {
        log.Print("\x1b[31;1m" + txt + "\x1b[0m")
}

