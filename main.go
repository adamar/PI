
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

    out, err := getProcUptime(*pid)
    if err != nil {
        log.Print(err)
    }
    log.Print("uptime")
    log.Print(out)

    out2, err := getEnv(*pid)
    if err != nil {
        log.Print(err)
    }
    log.Print("env")
    log.Print(out2)

    out3, err := getIO(*pid)
    if err != nil {
        log.Print(err)
    }
    log.Print("io")
    log.Print(out3)


}

func getProcUptime(pid string) (string, error) {

    comm := "ps"
    flags := []string{"-p", pid, "-o", "etime="}
    val, err := simpleRunCmd(comm, flags)
    if err != nil {
        return "", err
    }
    return val, nil

}


func getEnv(pid string) ([]string, error) {

    comm := "cat"
    log.Print(pid)
    flags := []string{"/proc/" + pid + "/environ"}
    val, err := runCmd(comm, flags)
    if err != nil {
        return nil, err
    }
    return val, nil

}


func getIO(pid string) ([]string, error) {

    comm := "cat"
    flags := []string{`/proc/` + pid + `/io`, `|`, `grep`,`"^bytes"`}
    val, err := runCmd(comm, flags)
    if err != nil {
        return nil, err
    }
    return val, nil

}


func simpleRunCmd(comm string, flags []string) (string, error) {

  cmd, err := exec.Command(comm, flags...).CombinedOutput()

  return string(cmd), err

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
        log.Print(line)
        log.Print(err)
        if err == io.EOF {
            break
        }

        output = append(output, line)

    }

    cmd.Wait()

    return output, nil

}



