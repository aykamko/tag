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
	"text/template"

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
)

type AliasFile struct {
	filename string
	fmtStr   string
	buf      bytes.Buffer
	writer   *bufio.Writer
}

func NewAliasFile() *AliasFile {
	aliasFilename := os.Getenv("TAG_ALIAS_FILE")
	if len(aliasFilename) == 0 {
		aliasFilename = "/tmp/tag_aliases"
	}

	aliasPrefix := os.Getenv("TAG_ALIAS_PREFIX")
	if len(aliasPrefix) == 0 {
		aliasPrefix = "e"
	}

	aliasCmdFmtString := os.Getenv("TAG_CMD_FMT_STRING")
	if len(aliasCmdFmtString) == 0 {
		aliasCmdFmtString = "vim {{.Filename}} +{{.LineNumber}}"
	}

	a := &AliasFile{
		fmtStr:   "alias " + aliasPrefix + "{{.MatchIndex}}='" + aliasCmdFmtString + "'\n",
		filename: aliasFilename,
	}
	a.writer = bufio.NewWriter(&a.buf)
	return a
}

func (a *AliasFile) WriteAlias(index int, filename, linenum string) {
	t := template.Must(template.New("alias").Parse(a.fmtStr))

	aliasVars := struct {
		MatchIndex int
		Filename   string
		LineNumber string
	}{index, filename, linenum}

	err := t.Execute(a.writer, aliasVars)
	check(err)
}

func (a *AliasFile) WriteFile() {
	err := a.writer.Flush()
	check(err)

	err = ioutil.WriteFile(a.filename, a.buf.Bytes(), 0644)
	check(err)
}

func optionIndex(args []string, option string) int {
	for i := len(args) - 1; i >= 0; i-- {
		if args[i] == option {
			return i
		}
	}
	return -1
}

func tagPrefix(aliasIndex int) string {
	return blue("[") + red("%d", aliasIndex) + blue("]")
}

func generateTags(cmd *exec.Cmd) int {
	color.NoColor = (optionIndex(cmd.Args, "--nocolor") > 0)
	cmd.Args = append(cmd.Args, "--null")

	stdout, err := cmd.StdoutPipe()
	check(err)
	scanner := bufio.NewScanner(stdout)

	var (
		line      []byte
		curPath   string
		groupIdxs []int
	)

	aliasFile := NewAliasFile()
	defer aliasFile.WriteFile()

	aliasIndex := 1

	err = cmd.Start()
	check(err)

	for scanner.Scan() {
		line = scanner.Bytes()
		if len(line) > 0 {
			if groupIdxs = pathRe.FindSubmatchIndex(line); len(groupIdxs) > 0 {
				// Extract path, print path, slice off path prefix
				curPath = string(line[groupIdxs[2]:groupIdxs[3]])
				fmt.Println(string(line[:groupIdxs[1]]))
				line = line[groupIdxs[1]:]
			}
			if groupIdxs = lineNumberRe.FindSubmatchIndex(line); len(groupIdxs) > 0 {
				aliasFile.WriteAlias(aliasIndex, curPath, string(line[groupIdxs[2]:groupIdxs[3]]))
				fmt.Printf("%s %s\n", tagPrefix(aliasIndex), string(line))
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

func isatty(f *os.File) bool {
	stat, err := f.Stat()
	check(err)
	return stat.Mode()&os.ModeCharDevice != 0
}

func main() {
	noTag := false
	var tagArgs []string

	switch i := optionIndex(os.Args, "--notag"); {
	case i > 0:
		os.Args = append(os.Args[:i], os.Args[i+1:]...)
		fallthrough
	case len(os.Args) == 1:
		noTag = true
	default:
		tagArgs = []string{"--group", "--color"}
	}

	/* From ag src/options.c:
	 * If we're not outputting to a terminal. change output to:
	 * turn off colors
	 * print filenames on every line
	 */
	args := os.Args[1:]
	if isatty(os.Stdout) {
		args = append(tagArgs, args...)
	}

	cmd := exec.Command("ag", args...)
	cmd.Stderr = os.Stderr

	if noTag || !isatty(os.Stdin) || !isatty(os.Stdout) {
		// Data being piped from stdin
		os.Exit(passThrough(cmd))
	}

	os.Exit(generateTags(cmd))
}
