
package main

import (
	"bytes"
        "os"
	"os/exec"
        "io"
        "io/ioutil"
        "flag"
        "log"
        "strings"
        "errors"
)



func main() {

    var pid = flag.String("p", "", "pid of the process to inspect")
    flag.Parse()

    if *pid == "" {
        os.Exit(1)
    }


    out, err := getProcUptime(*pid)
    if err != nil {
        log.Print("uptime error")
    }
    log.Print("uptime")
    log.Print(out)

    //out2, err := getEnv(*pid)
    //if err != nil {
    //    log.Print("env error")
    //}
    //log.Print("env")
    //log.Print(out2)

    out3, err := getIO(*pid)
    if err != nil {
        log.Print("io error")
    }
    log.Print("io")
    log.Print(out3)

    out4, err := getProcStatus(*pid)
    if err != nil {
        log.Print("proc stat error")
    }
    log.Print("state")
    log.Print(out4)



}

func getProcStatus(pid string) (string, error) {

    val, err := ioutil.ReadFile("/proc/" + pid + "/stat")
    if err != nil {
        return "", err
    }

    output := strings.Split(string(val), " ")

    var value string

    switch output[2] {
            case "Z":
              value = "Zombie"
            case "S":
              value = "Sleeping"
            case "R":
              value = "Running"
            case "D":
              value = "Waiting in Disk Sleep"
            case "W":
              value = "Paging"
    }
    return value, nil

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

    var fileName = "/proc/" + pid + "/environ"

    fileExists, _ := fileExists(fileName)
    if fileExists == false {
        return nil, errors.New("File doesnt exist")
    }

    val, err := ioutil.ReadFile(fileName)
    if err != nil {
        return nil, err
    }

    output := strings.Split(string(val), "\000")
    return output, nil

}


func getIO(pid string) (string, error) {

    val, err := ioutil.ReadFile("/proc/" + pid + "/io")
    if err != nil {
        return "", err
    }

    return string(val), nil

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
        //log.Print(line)
        //log.Print(err)
        if err == io.EOF {
            break
        }

        output = append(output, line)

    }

    cmd.Wait()

    return output, nil

}

func fileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { 
        return true, nil 
    }
    if os.IsNotExist(err) { 
        return false, nil 
    }
    return false, err
}
