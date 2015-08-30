// This provides a glog style logging for web logs. In fact parts of thh code is copied from glog source.
package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}
}

func createLogDir(name string) error {
	err := os.MkdirAll(name, 0775)
	if err != nil {
		return err
	}
	return nil
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func logName(l string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%s.%s.%s.%04d%02d%02d-%02d%02d%02d.%d",
		program,
		host,
		userName,
		l,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
	return name, program + "." + l
}

func setupWebLog(dirPath string, t time.Time) (fi *os.File, err error) {
	e := createLogDir(dirPath)
	if e != nil {
		return nil, fmt.Errorf("log: cannot create log: %v", e)
	}
	name, link := logName("access.log", t)
	fname := filepath.Join(dirPath, name)
	f, err := os.Create(fname)
	if err == nil {
		symlink := filepath.Join(dirPath, link)
		os.Remove(symlink)        // ignore err
		os.Symlink(name, symlink) // ignore err
		return f, nil
	}
	return nil, fmt.Errorf("log: cannot create log: %v", err)
}
