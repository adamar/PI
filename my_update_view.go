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
