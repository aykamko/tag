package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/fatih/color"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var (
	red          = color.RedString
	blue         = color.BlueString
	pathRe       = regexp.MustCompile(`^(?:\x1b\[[^m]+m)?([^\x00\x1b]+)[^\x00]*\x00`)
	lineNumberRe = regexp.MustCompile(`^(?:\x1b\[[^m]+m)?(\d+)(?:\x1b\[0m\x1b\[K)?:`)
	cwd          string
)

type AliasFile struct {
	buf    bytes.Buffer
	writer *bufio.Writer
}

func NewAliasFile() *AliasFile {
	a := &AliasFile{}
	a.writer = bufio.NewWriter(&a.buf)
	return a
}

func (a *AliasFile) WriteAlias(index int, filename, linenum string) {
	fmt.Fprintf(a.writer, "alias f%d='vim %s/%s +%s'\n", index, cwd, filename, linenum)
}

func (a *AliasFile) WriteFile(filename string) {
	err := a.writer.Flush()
	check(err)

	err = ioutil.WriteFile(filename, a.buf.Bytes(), 0644)
	check(err)
}

func hasOption(args []string, option string) bool {
	for i := len(args) - 1; i >= 0; i-- {
		if args[i] == option {
			return true
		}
	}
	return false
}

func generateShortcuts(cmd *exec.Cmd) int {
	color.NoColor = hasOption(cmd.Args, "--nocolor")
	cmd.Args = append(cmd.Args, "--null")

	stdout, err := cmd.StdoutPipe()
	check(err)
	scanner := bufio.NewScanner(stdout)

	cwd, err = os.Getwd()
	check(err)

	var (
		line      []byte
		curPath   string
		groupIdxs []int
	)
	aliasFile := NewAliasFile()
	defer aliasFile.WriteFile("/tmp/tag_aliases")
	aliasIndex := 1

	err = cmd.Start()
	check(err)

	for scanner.Scan() {
		line = scanner.Bytes()
		if len(line) > 0 {
			if groupIdxs = pathRe.FindSubmatchIndex(line); len(groupIdxs) > 0 {
				// Extract path, print path, trancate path prefix
				curPath = string(line[groupIdxs[2]:groupIdxs[3]])
				fmt.Println(string(line[:groupIdxs[1]]))
				line = line[groupIdxs[1]:]
			}
			if groupIdxs = lineNumberRe.FindSubmatchIndex(line); len(groupIdxs) > 0 {
				aliasFile.WriteAlias(aliasIndex, curPath, string(line[groupIdxs[2]:groupIdxs[3]]))
				fmt.Printf("%s%s%s %s\n", blue("["), red("%d", aliasIndex), blue("]"), string(line))
				aliasIndex++
				continue
			}
		}
		fmt.Println(string(line))
	}

	err = cmd.Wait()
	if err != nil {
		return 1
	}
	return 0
}

func passThrough(cmd *exec.Cmd) int {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return 1
	}
	return 0
}

func main() {
	args := []string{"--group", "--color"}
	args = append(args, os.Args[1:]...)

	cmd := exec.Command("ag", args...)
	cmd.Stderr = os.Stderr

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 { // Data being piped from stdin
		os.Exit(passThrough(cmd))
	}

	os.Exit(generateShortcuts(cmd))
}
