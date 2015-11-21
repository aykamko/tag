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
	red        = color.RedString
	blue       = color.BlueString
	ansi       = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})*)?[a-zA-Z]`)
	lineNumber = regexp.MustCompile(`^(\d+):`)
	cwd        string
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

func main() {
	cmd := exec.Command("ag", append([]string{"--group", "--color"}, os.Args[1:]...)...)
	stdout, err := cmd.StdoutPipe()
	check(err)

	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	check(err)

	cwd, err = os.Getwd()
	check(err)

	var (
		line         string
		strippedLine string
		curfile      string
		groups       []string
	)
	aliasFile := NewAliasFile()
	aliasIndex := 1

	for scanner.Scan() {
		line = scanner.Text()
		strippedLine = stripAnsi(line)
		if len(line) > 0 {
			groups = lineNumber.FindStringSubmatch(strippedLine)
			if len(groups) > 0 {
				aliasFile.WriteAlias(aliasIndex, curfile, stripAnsi(groups[1]))
				fmt.Printf("%s%s%s %s\n", blue("["), red(fmt.Sprintf("%d", aliasIndex)), blue("]"), line)
				aliasIndex++
				continue
			}

			curfile = strippedLine
		}
		fmt.Println(line)
	}

	aliasFile.WriteFile("/tmp/tag_aliases")
}

func stripAnsi(str string) string {
	return ansi.ReplaceAllString(str, "")
}
