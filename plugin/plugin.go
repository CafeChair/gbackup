package plugin

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Script struct {
	ScriptName string
	Cycle      int
}

func ListScripts(relativepath string) []*Script {
	// ret := make(map[string]*Script)
	// script := new(Script)
	scripts := make([]*Script, 0)
	fs, err := ioutil.ReadDir(relativepath)
	if err != nil {
		log.Println("can't list file under: ", relativepath)
		// return ret
		return scripts
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		arr := strings.Split(filename, "_")
		if len(arr) < 2 {
			continue
		}
		var cycle int
		cycle, err = strconv.Atoi(arr[0])
		if err != nil {
			continue
		}
		scriptname := filepath.Join(relativepath, filename)
		// scripts := &Script{ScriptName: scriptname, Cycle: cycle}
		script := &Script{ScriptName: scriptname, Cycle: cycle}
		// ret[scriptname] = scripts
		scripts = append(scripts, script)
	}
	// return ret
	return scripts
}

func RunScript(script *Script) (string, error) {
	timeout := script.Cycle * 1000
	cmd := exec.Command(script.ScriptName)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// err := cmd.Run()
	// return stdout.String(), err
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Start()
	err, IsTimeout := scriptRunTimeOut(cmd, time.Duration(timeout)*time.Millisecond)
	errStr := stderr.String()
	if errStr != "" {
		log.Println("[Error] exec script", script.ScriptName, "fail. error: ", errStr)
		// return "", err
	}
	if IsTimeout {
		if err == nil {
			log.Println("[Info] timeout and kill process", script.ScriptName, "success")
		}
		if err != nil {
			log.Println("[Error] timeout and kill process", script.ScriptName, "error:", err)
		}
		return "", err
		// return
	}
	return stdout.String(), err
}

func scriptRunTimeOut(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	var err error
	select {
	case <-time.After(timeout):
		log.Printf("timeout, process:%s will be killed", cmd.Path)
		go func() {
			<-done
		}()
		if err = cmd.Process.Kill(); err != nil {
			log.Printf("fail to kill: %s, error: %s", cmd.Path, err)
		}
		return err, true
	case err = <-done:
		return err, false
	}
}
