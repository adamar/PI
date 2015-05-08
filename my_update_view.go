package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os/exec"
	//"sort"
	"strings"
	"time"
)

// PI is the main display struct
type PI struct {
	Syscalls map[string]int
	Current  string
	GUI      *gocui.Gui
	Pid      string
}

func main() {

	pi := NewPI()
	pi.Start()

}

// NewPI returns a new PI instance
func NewPI() *PI {

	pi := PI{}

	var err error
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		panic(err)
	}

	g.SetLayout(pi.layout)
	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		panic(err)
	}); err != nil {
		log.Panicln(err)
	}

	pi.Syscalls = map[string]int{}
	pi.GUI = g
	pi.Pid = "14295"
	return &pi

}

func (pi *PI) Start() {
	defer pi.GUI.Close()

	go pi.run()
	go pi.updateView()
	err := pi.GUI.MainLoop()
	if err != nil {
		panic(err)
	}

}

func (pi *PI) run() {

	ch := make(chan string)

	go func(ch chan string) {
		cmd := exec.Command("sudo", "strace", "-f", "-e", "trace=network", "-p", pi.Pid)
		stdout, _ := cmd.StderrPipe()
		cmd.Start()
		r := bufio.NewReader(stdout)

		for {

			bufline, _ := r.ReadString('\n')
			parsedSyscall := parseSyscallString(bufline)
			ch <- parsedSyscall

		}

	}(ch)

	for {
		data := <-ch

		//if HasPrefix(data, "Attached") <--- Ignore Process n attached

		if _, ok := pi.Syscalls[data]; ok {
			pi.Syscalls[data] = pi.Syscalls[data] + 1
		} else {
			pi.Syscalls[data] = 1
		}
		//pi.Syscalls = append(pi.Syscalls, data)
		pi.Current = data
	}

}

func (pi *PI) updateView() {
	for {
		time.Sleep(time.Second)
		v, err := pi.GUI.View("center")
		if err != nil {
			panic(err)
		}
		v.Clear()
		pi.display(v)
		pi.GUI.Flush()
	}
}

func (pi *PI) display(v io.Writer) error {

	output := ""

	//sort.Strings(pi.Syscalls)

	for k, _ := range pi.Syscalls {
		if k == pi.Current {
			output += fmt.Sprintf("=>" + k + "\n")
		} else {
			output += fmt.Sprintf(k + "\n")
		}
	}

	fmt.Fprintf(v, output)
	return nil
}

func (pi *PI) layout(*gocui.Gui) error {
	maxX, maxY := pi.GUI.Size()
	if v, err := pi.GUI.SetView("center", 3, 0, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Frame = false
		pi.display(v)
	}
	return nil
}

func parseSyscallString(bufline string) string {

	parsed, _ := parseAlphanumeric(bufline)

	return strings.Join(parsed[:4], " ")

}

func parseAlphanumeric(str string) ([]string, bool) {
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
