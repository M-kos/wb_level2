package mycut

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type CutOptions struct {
	Fields    []int
	Delimiter string
	Separate  bool
}

var (
	ErrFieldsFlag   = errors.New("-f is required")
	ErrInvalidRange = errors.New("invalid range")
	ErrInvalidField = errors.New("invalid field")
	ErrReadInput    = errors.New("error reading input")
)

func Run() error {
	options, err := parseFleg()
	if err != nil {
		return err
	}

	res, err := runCut(os.Stdin, options)
	if err != nil {
		return err
	}

	for _, line := range res {
		fmt.Println(line)
	}

	return nil
}

func runCut(reader io.Reader, options *CutOptions) ([]string, error) {
	scanner := bufio.NewScanner(reader)

	result := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if options.Separate && !strings.Contains(line, options.Delimiter) {
			continue
		}

		parts := strings.Split(line, options.Delimiter)
		selected := make([]string, 0, len(options.Fields))

		for _, idx := range options.Fields {
			if idx < len(parts) {
				selected = append(selected, parts[idx])
			}
		}

		result = append(result, strings.Join(selected, options.Delimiter))
	}

	if err := scanner.Err(); err != nil {
		return nil, ErrReadInput
	}

	return result, nil
}

func parseFleg() (*CutOptions, error) {
	fStr := flag.String("f", "", "fields")
	delim := flag.String("d", "\t", "delimiter")
	separated := flag.Bool("s", false, "only lines with delimiter")
	flag.Parse()

	if *fStr == "" {
		return nil, ErrFieldsFlag
	}

	fields, err := parseFields(*fStr)
	if err != nil {
		return nil, err
	}

	return &CutOptions{
		Fields:    fields,
		Delimiter: *delim,
		Separate:  *separated,
	}, nil
}

func parseFields(fieldsStr string) ([]int, error) {
	fields := make([]int, 0)
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")

			if len(bounds) != 2 {
				return nil, ErrInvalidRange
			}

			start, errStart := strconv.Atoi(bounds[0])
			end, errEnd := strconv.Atoi(bounds[1])

			if errStart != nil || errEnd != nil || start < 1 || end < start {
				return nil, ErrInvalidRange
			}

			for i := start - 1; i < end; i++ {
				fields = append(fields, i)
			}

			continue
		}

		idx, err := strconv.Atoi(part)
		if err != nil || idx < 1 {
			return nil, ErrInvalidField
		}

		fields = append(fields, idx-1)
	}

	slices.Sort(fields)

	return fields, nil
}
