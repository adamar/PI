
package main

import (
	"bytes"
        "os"
	"os/exec"
        "io"
        "flag"
        "log"
)


func main() {

    var pid = flag.String("p", "", "pid of the process to inspect")
    flag.Parse()

    if *pid == "" {
        os.Exit(1)
    }


    comm := "ps"
    flags := []string{"-p", *pid, "-o", "etime="}
    val, _ := runCmd(comm, flags)
    log.Print(val)

}


func getEnv(pid string) error {

    comm := "cat"
    flags := []string{"/proc/" + *pid + "/environ"}
    val, _ := runCmd(comm, flags)
    return nil

}


func getIO(pid string) error {

    comm := "cat"
    flags := []string{`/proc/` + *pid + `/io`, `|`, `grep`,`"^bytes"`}
    val, _ := runCmd(comm, flags)
    return nil

}





func runCmd(comm string, flags []string) ([]string, error) {

    cmd := exec.Command(comm, flags...)

    output := []string{}

    // STDOUT
    stdPipe, err := cmd.StdoutPipe()
    if err != nil {
        os.Exit(0)
        return nil, err
    }

    // STDERR
    errPipe, err := cmd.StderrPipe()
    if err != nil {
        os.Exit(0)
        return nil, err
    }

    // Exec the command
    err = cmd.Start()
    if err != nil {
        os.Exit(0)
        return nil, err
    }

    errBuf := new(bytes.Buffer)
    errBuf.ReadFrom(errPipe)

    io.Copy(os.Stdout, stdPipe)
    stdBuf := new(bytes.Buffer)
    stdBuf.ReadFrom(stdPipe)

    for {

        line, err := stdBuf.ReadString('\n')
        if err == io.EOF {
            break
        }

        output = append(output, line)

    }

    cmd.Wait()

    return output, nil

}



