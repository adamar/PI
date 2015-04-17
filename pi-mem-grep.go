package main

import (
	"log"
	"os"
	"os/exec"
)

type Command struct {
	Name string
	Args []string
	Grep string
}

type Result struct {
	Output []string
}

func (c *Command) execCommand() error {

	bin, err := exec.LookPath(c.Name)
	if err != nil {
		return err
	}

	output, err := exec.Command(bin, c.Args...).CombinedOutput()
	if err != nil {
		return err
	}

	log.Print(string(output))

	return nil

}

func checkPid(pid string) error {

	_, err := os.Stat("/proc/" + pid + "/maps")
	if err != nil {
		return err
	}

	return nil

}

func buildCommand(command []string, grep string) *Command {
	return &Command{Name: command[0], Args: command[1:], Grep: grep}
}

func main() {

	arr := []string{"gdp", "--batch", "-pid", "30182", "-ex", `"dump memory 0x7fff145c4000 0x7fff145e5000"`}
	ex := buildCommand(arr, "grepstring")
	ex.execCommand()

}
