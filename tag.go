package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	red  = color.RedString
	blue = color.BlueString
	ansi = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})*)?[a-zA-Z]`)
	cwd  string
)

func stripAnsi(str string) string {
	return ansi.ReplaceAllString(str, "")
}

func writeAlias(writer io.Writer, index int, filename, linenum string) {
	fmt.Fprintf(writer, "alias f%d='vim %s/%s +%s'\n", index, cwd, filename, linenum)
}

func main() {
	cmd := exec.Command("ag", "--group", "--color", os.Args[1])
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var (
		aliasBuf   bytes.Buffer
		line       string
		curfile    string
		linenumEnd int
		aliasIndex int
	)
	bufWriter := bufio.NewWriter(&aliasBuf)
	aliasIndex = 1
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > 0 {
			linenumEnd = strings.IndexByte(line, ':')
			if linenumEnd > 0 {
				writeAlias(bufWriter, aliasIndex, curfile, stripAnsi(line[:linenumEnd]))
				fmt.Printf("%s%s%s %s\n", blue("["), red(fmt.Sprintf("%d", aliasIndex)), blue("]"), line)
				aliasIndex++
				continue
			}

			curfile = stripAnsi(line)
		}
		fmt.Println(line)
	}

	err = bufWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("/tmp/tag_aliases", aliasBuf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
