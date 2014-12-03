// Package pstree defines processes and process trees, with the ability to populate them
package pstree

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// ProcessIDs are strings. We don't do math and creation/consumption are guarded with digit globs and digit regexes.
type ProcessID string

// A Process knows itself, its parent, and who its children are.
type Process struct {
	name     string
	pid      ProcessID
	parent   ProcessID
	children []ProcessID
}

// ProcessTree represents every process on a system.
type ProcessTree map[ProcessID]*Process

// ReadProcessInfo parses the /proc/pid/stat file to fill the struct.
// /proc/pid/stat does not contain information about children -- only parents.
// Linking children to parents is done in a later pass when we know about every process.
func (proc *Process) ReadProcessInfo(pid ProcessID) (err error) {

	filename := "/proc/" + string(pid) + "/stat"

	fh, err := os.Open(filename)
	if err != nil {
		return err
	}

	// 512 being a safe/sloppy upper bound for the size of a stat file
	var line = make([]byte, 512)
	if _, err := fh.Read(line); err != io.EOF && err != nil {
		return err
	}
	fh.Close()

	// 25926 (annoy me.out) S 25906 31864 31842 ...
	// 2nd entry is the name in parens, 4th is parent pid
	// get the index of the first and last paren to grab process name
	first := bytes.IndexByte(line, '(')
	last := bytes.IndexByte(line[first:], ')') + first
	if first == -1 || last == -1 {
		return errors.New("Can't parse " + filename)
	}

	// don't take the parens with us
	proc.name = string(line[first+1 : last])

	// skip ahead to the rest of the string starting at PPID
	rest := line[last+4:]
	last = bytes.IndexByte(rest, ' ')
	if last == -1 {
		return errors.New("Can't parse " + filename)
	}
	proc.parent = ProcessID(rest[:last])
	proc.pid = pid

	return nil
}

// Populate reads the entries in /proc and then assembles the list of children.
// In order to link parents to children, we must first know every processes's parent.
func (processes ProcessTree) Populate() error {

	// create pid 0 manually since /proc doesn't expose it, but processes have parents that are pid 0
	processes["0"] = &Process{"sched", "0", "", nil}

	procfh, err := os.Open("/proc/")
	if err != nil {
		return err
	}

	entries, err := procfh.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !isInt(entry) {
			continue
		}

		pid := ProcessID(entry)
		processes[pid] = new(Process)
		if err := processes[pid].ReadProcessInfo(pid); err != nil {
			return err
		}

	}

	// for every process, add itself to its parent's list of children, now that we know parents
	for pid, info := range processes {
		if pid == "0" {
			continue
		}
		processes[info.parent].children = append(processes[info.parent].children, pid)
	}

	return nil
}

// PrintDepthFirst traverses the hash and recursively prints a parent's children.
func (pids ProcessTree) PrintDepthFirst(pid ProcessID, depth int) string {
	res := fmt.Sprintf("%*s%v (%v)\n", depth, "", pids[pid].name, pid)
	for _, kid := range pids[pid].children {
		res = res + pids.PrintDepthFirst(kid, depth+1)
	}
	return res
}

// String assumes the user wants all processes and is a print-friendly wrapper of PrintDepthFirst
func (pids ProcessTree) String() string {
	return pids.PrintDepthFirst("0", 0)
}

func isInt(in string) bool {
	for i := 0; i < len(in); i++ {
		var b byte = in[i]
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
