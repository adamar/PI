
package main

import (
        "os/exec"
        "bufio"
        "log"
        //"path/filepath"
        //"io"
        "strings"
        "flag"
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




}




