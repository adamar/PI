
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

  for {

      bufline, err := r.ReadString('\n')
      line := string(bufline)

      lines := strings.Split(line, " ")
      filestring := strings.Split(lines[3], `"`)

      switch {

          case strings.HasPrefix(lines[3], "open"):
              log.Print("open ", filestring[1])

          case strings.HasPrefix(lines[3], "stat"):
              log.Print("stat ", filestring[1])

          case strings.HasPrefix(lines[3], "readlink"):
              log.Print("check symlink ", filestring[1])



          default:
             log.Print(lines[3])

      }

      if err != nil {
         log.Fatal(err)
      }
  }



}




