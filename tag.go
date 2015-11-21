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
	red  = color.RedString
	blue = color.BlueString
	cwd  string
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

func getOutputColors(args []string) (pathColor string, lineNumberColor string) {
	pathColor, lineNumberColor = "1;32", "1;33"
	for i, flag := range args {
		if flag == "--color-path" && i+1 < len(args) {
			pathColor = args[i+1]
		} else if flag == "--color-line-number" && i+1 < len(args) {
			lineNumberColor = args[i+1]
		}
	}
	return
}

func generateShortcuts(cmd *exec.Cmd) {
	pathColor, lineNumberColor := getOutputColors(os.Args[1:])
	pathRe := regexp.MustCompile(
		fmt.Sprintf(`^\x1b\[%sm([^\x1b]+)`, pathColor))
	lineNumberRe := regexp.MustCompile(
		fmt.Sprintf(`^\x1b\[%sm(\d+)\x1b\[0m\x1b\[K:`, lineNumberColor))

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
	aliasIndex := 1

	err = cmd.Start()
	check(err)

	for scanner.Scan() {
		line = scanner.Bytes()
		if len(line) > 0 {
			if groupIdxs = lineNumberRe.FindSubmatchIndex(line); len(groupIdxs) > 0 {
				aliasFile.WriteAlias(aliasIndex, curPath, string(line[groupIdxs[2]:groupIdxs[3]]))
				fmt.Printf("%s%s%s %s\n", blue("["), red("%d", aliasIndex), blue("]"), string(line))
				aliasIndex++
				continue
			} else if groupIdxs = pathRe.FindSubmatchIndex(line); len(groupIdxs) > 0 {
				curPath = string(line[groupIdxs[2]:groupIdxs[3]])
			}
		}
		fmt.Println(string(line))
	}

	aliasFile.WriteFile("/tmp/tag_aliases")
}

func passThrough(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}

func main() {
	args := []string{"--group"}
	args = append(args, os.Args[1:]...)
	args = append(args, "--color")

	cmd := exec.Command("ag", args...)

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 { // Data being piped from stdin
		passThrough(cmd)
		return
	}

	generateShortcuts(cmd)
}
