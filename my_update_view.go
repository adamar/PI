package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os/exec"
	//"sort"
	"time"
)

// PI is the main display struct
type PI struct {
	Syscalls map[string]int
	Current  string
	GUI      *gocui.Gui
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
