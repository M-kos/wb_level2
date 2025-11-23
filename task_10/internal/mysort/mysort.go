package mysort

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	Raw       string
	Key       string
	Numeric   float64
	Month     int
	HumanSize int64
}

type SortOptions struct {
	Column              uint
	ReverseSort         bool
	NumericSort         bool
	Unique              bool
	MonthSort           bool
	IgnoreLeadingBlanks bool
	CheckSorted         bool
	HumanNumericSort    bool
}

var monthMap = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
	"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
	"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

var (
	ErrInvalidArgumentForTheKflag = errors.New("invalid argument for the -k flag")
	ErrNeedPathToFile             = errors.New("need path to file")
)

func Run() ([]string, error) {
	if len(os.Args) < 2 {
		return nil, ErrNeedPathToFile
	}

	filePath := getFilePath()
	reader := getReader(filePath)
	return runSort(reader)
}

func runSort(reader io.Reader) ([]string, error) {
	opt, err := parseFlags()
	if err != nil {
		return nil, err
	}

	lines, err := readLines(reader, opt)
	if err != nil {
		return nil, err
	}

	if opt.CheckSorted {
		if checkSorted(lines, opt) {
			fmt.Println("sorted")
		} else {
			fmt.Println("unsorted")
		}
		return nil, nil
	}

	lines = sortLines(lines, opt)

	result := make([]string, 0, len(lines))
	for _, l := range lines {
		result = append(result, l.Raw)
	}

	return result, nil
}

func checkSorted(lines []Line, options *SortOptions) bool {
	return sort.SliceIsSorted(lines, func(i, j int) bool {
		a := lines[i]
		b := lines[j]

		switch {
		case options.MonthSort:
			if a.Month != b.Month {
				return a.Month < b.Month
			}
		case options.HumanNumericSort:
			if a.HumanSize != b.HumanSize {
				return a.HumanSize < b.HumanSize
			}
		case options.NumericSort:
			if a.Numeric != b.Numeric {
				return a.Numeric < b.Numeric
			}
		}
		return a.Key < b.Key
	})
}

func sortLines(lines []Line, options *SortOptions) []Line {
	sort.Slice(lines, func(i, j int) bool {
		a := lines[i]
		b := lines[j]

		if options.ReverseSort {
			a, b = b, a
		}

		switch {
		case options.MonthSort:
			if a.Month != b.Month {
				return a.Month < b.Month
			}
		case options.HumanNumericSort:
			if a.HumanSize != b.HumanSize {
				return a.HumanSize < b.HumanSize
			}
		case options.NumericSort:
			if a.Numeric != b.Numeric {
				return a.Numeric < b.Numeric
			}
		}
		return a.Key < b.Key
	})

	if options.Unique {
		unique := make([]Line, 0, len(lines))
		var prev string
		for _, l := range lines {
			if l.Key != prev {
				unique = append(unique, l)
				prev = l.Key
			}
		}
		lines = unique
	}

	return lines
}

func makeLine(str string, options *SortOptions) Line {
	key := getColumn(str, options.Column)
	if key == "" {
		key = str
	}

	if options.IgnoreLeadingBlanks {
		key = strings.TrimLeft(key, " ")
	}

	line := Line{
		Raw: str,
		Key: key,
	}

	if options.MonthSort {
		if val, ok := monthMap[key]; ok {
			line.Month = val
		}
	}

	if options.NumericSort {
		line.Numeric, _ = strconv.ParseFloat(key, 64)
	}

	if options.HumanNumericSort {
		line.HumanSize, _ = getHumanSize(key)
	}

	return line
}

func getHumanSize(s string) (int64, error) {
	mult := int64(1)
	if len(s) == 0 {
		return 0, errors.New("empty")
	}

	last := s[len(s)-1]
	switch last {
	case 'K', 'k':
		mult = 1024
		s = s[:len(s)-1]
	case 'M', 'm':
		mult = 1024 * 1024
		s = s[:len(s)-1]
	case 'G', 'g':
		mult = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case 'T', 't':
		mult = 1024 * 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return int64(num * float64(mult)), nil
}

func getColumn(line string, columnNumber uint) string {
	if columnNumber == 0 {
		return ""
	}

	columns := strings.Split(line, "\t")
	columnIdx := int(columnNumber) - 1
	if columnIdx < len(columns) {
		return columns[columnIdx]
	}

	return ""
}

func readLines(reader io.Reader, options *SortOptions) ([]Line, error) {
	sc := bufio.NewScanner(reader)
	var lines []Line

	for sc.Scan() {
		lines = append(lines, makeLine(sc.Text(), options))
	}

	return lines, sc.Err()
}

func getReader(filePath string) io.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		return os.Stdin
	}

	return file
}

func parseFlags() (*SortOptions, error) {
	options := &SortOptions{}
	flags := make([]string, 0)

	for _, arg := range os.Args[1:] {
		_, err := os.Stat(arg)
		if err == nil {
			continue
		}

		f := strings.TrimLeft(arg, "-")
		flags = append(flags, f)
	}

	flagsStr := strings.Join(flags, "")

	for i := 0; i < len(flagsStr); {
		switch flagsStr[i] {
		case 'k':
			if i+1 >= len(flagsStr) {
				return nil, ErrInvalidArgumentForTheKflag
			}

			n, err := strconv.Atoi(string(flagsStr[i+1]))
			if err != nil || n <= 0 {
				return nil, ErrInvalidArgumentForTheKflag
			}

			options.Column = uint(n)
			i = i + 2
			continue
		case 'r':
			options.ReverseSort = true
		case 'n':
			options.NumericSort = true
		case 'u':
			options.Unique = true
		case 'M':
			options.MonthSort = true
		case 'b':
			options.IgnoreLeadingBlanks = true
		case 'c':
			options.CheckSorted = true
		case 'h':
			options.HumanNumericSort = true
		default:
			return nil, fmt.Errorf("unknown flag: -%c", flagsStr[i])
		}
		i++
	}

	return options, nil
}

func getFilePath() string {
	if len(os.Args) < 2 {
		return ""
	}

	for _, arg := range os.Args[1:] {
		_, err := os.Stat(arg)
		if err == nil {
			return arg
		}
	}

	return ""
}
