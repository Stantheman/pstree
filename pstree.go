package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// PIDs are strings -- we don't do math and creation/consumption are guarded with digit globs and digit regexes.
type ProcessID string

// A Process knows itself, its parent, and who its children are.
type Process struct {
	name     string
	pid      ProcessID
	parent   ProcessID
	children []ProcessID
}

// ProcessTree is just a convenient world-view of all processes
type ProcessTree map[ProcessID]*Process

const IndentDepth = 0

func main() {
	tree := make(ProcessTree)

	if err := tree.Populate(); err != nil {
		fmt.Errorf("Failed getting proceses: %v\n", err)
		return
	}

	tree.PrintDepthFirst("0", IndentDepth)
}

func (processes ProcessTree) Populate() error {
	matches, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return err
	}

	for _, pidpath := range matches {
		pid := ProcessID(filepath.Base(pidpath))

		processes[pid] = new(Process)
		if err := processes[pid].ReadProcessInfo(pid); err != nil {
			return err
		}
	}

	// cheat and create pid 0 since /proc doesn't expose it, but processes have parents that are pid 0
	processes["0"] = &Process{"sched", "0", "", nil}

	// now that we have the list of pids, populate the child list
	for pid, info := range processes {
		if pid == "0" {
			continue
		}
		processes[info.parent].children = append(processes[info.parent].children, pid)
	}

	return nil
}

func (pids ProcessTree) PrintDepthFirst(pid ProcessID, depth int) {
	fmt.Printf("%*s%v (%v)\n", depth, "", pids[pid].name, pid)
	for _, kid := range pids[pid].children {
		pids.PrintDepthFirst(kid, depth+1)
	}
}

func (proc *Process) ReadProcessInfo(pid ProcessID) (err error) {
	fh, err := os.Open("/proc/" + string(pid) + "/stat")
	if err != nil {
		return err
	}
	defer fh.Close()

	buf := bufio.NewReader(fh)
	// /proc/pid/stat is one line
	line, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	// 25926 (a.out) S 25906 31864 31842 ...
	// 2nd entry is the name in parens, 4th is parent pid
	re, err := regexp.Compile(`\(([\w\s\/\.:-]+)\)\s[A-Z]\s(\d+)`)
	if err != nil {
		return err
	}

	matches := re.FindStringSubmatch(line)
	if len(matches) < 3 {
		return errors.New("Couldn't match on: " + line)
	}

	proc.name = matches[1]
	proc.parent = ProcessID(matches[2])

	return nil
}
