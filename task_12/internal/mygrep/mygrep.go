package mygrep

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type GrepOptions struct {
	ContextAfter    int
	ContextBefore   int
	ContextAround   int
	Count           bool
	Ignore          bool
	Invert          bool
	Fixed           bool
	ShowLinesNumber bool
	Pattern         string
}

var (
	ErrNeedPattern = errors.New("need pattern")
	ErrOpenFile    = errors.New("error opening file")
)

func Run() error {
	if len(os.Args) < 2 {
		return ErrNeedPattern
	}

	reader := getReader()
	defer func() {
		err := reader.Close()
		if err != nil {
			fmt.Println("error closing:", err)
		}
	}()

	return runGrep(reader)
}

func runGrep(reader io.Reader) error {
	opt, err := parseFlags()
	fmt.Println(opt)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(reader)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	var lineMatches []int

	for i, line := range lines {
		match := findMatch(opt, line)
		if opt.Invert {
			match = !match
		}

		if match {
			lineMatches = append(lineMatches, i)
		}
	}

	if opt.Count {
		fmt.Println(len(lineMatches))
		return nil
	}

	for _, lineIdx := range lineMatches {
		for i := lineIdx - opt.ContextBefore; i < lineIdx; i++ {
			if i >= 0 {
				printLine(opt, lines[i], i+1)
			}
		}

		printLine(opt, lines[lineIdx], lineIdx+1)

		for i := lineIdx + 1; i <= lineIdx+opt.ContextAfter && i < len(lines); i++ {
			printLine(opt, lines[i], i+1)
		}
	}

	return nil
}

func parseFlags() (*GrepOptions, error) {
	after := flag.Int("A", 0, "after each line found, print an additional N lines after it.")
	before := flag.Int("B", 0, "print N lines before to each found line")
	around := flag.Int("C", 0, "print N lines around the found string")
	count := flag.Bool("c", false, "count of matching lines")
	ignore := flag.Bool("i", false, "ignore case")
	invert := flag.Bool("v", false, "invert the filter: output lines that do not contain a template")
	fixed := flag.Bool("F", false, "treat the template as a fixed string")
	showLineNum := flag.Bool("n", false, "show line numbers")

	flag.Parse()

	if *around > 0 {
		*after = *around
		*before = *around
	}

	if flag.NArg() < 1 {
		return nil, ErrNeedPattern
	}

	pattern := ""

	args := flag.Args()

	for _, a := range args {
		if _, err := os.Stat(a); err != nil {
			pattern = a
		}
	}

	options := &GrepOptions{
		ContextAfter:    *after,
		ContextBefore:   *before,
		ContextAround:   *around,
		Count:           *count,
		Ignore:          *ignore,
		Invert:          *invert,
		Fixed:           *fixed,
		ShowLinesNumber: *showLineNum,
		Pattern:         pattern,
	}

	return options, nil
}

func getReader() io.ReadCloser {
	for _, arg := range os.Args[1:] {
		_, err := os.Stat(arg)
		if err == nil {
			file, err := os.Open(arg)
			if err != nil {
				fmt.Println(errors.Join(ErrOpenFile, err))
				os.Exit(1)
			}

			return file
		}
	}

	return os.Stdin
}

func findMatch(options *GrepOptions, targetStr string) bool {
	if options.Fixed {
		if options.Ignore {
			return strings.Contains(strings.ToLower(targetStr), strings.ToLower(options.Pattern))
		}

		return strings.Contains(targetStr, options.Pattern)
	}

	pattern := options.Pattern

	if options.Ignore {
		pattern = "(?i)" + pattern
	}

	re := regexp.MustCompile(pattern)

	return re.MatchString(targetStr)
}

func printLine(options *GrepOptions, targetLine string, lineNum int) {
	if options.ShowLinesNumber {
		fmt.Printf("%d:%s\n", lineNum, targetLine)
	} else {
		fmt.Println(targetLine)
	}
}
