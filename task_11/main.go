package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println(anagrams([]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}))
}

func anagrams(input []string) map[string][]string {
	temp := make(map[string]string)
	result := make(map[string][]string)

	for _, v := range input {
		ordered := sortChars(v)
		lowV := strings.ToLower(v)

		if str, ok := temp[ordered]; ok {
			if _, ok := result[str]; ok {
				result[str] = append(result[str], lowV)
				continue
			}

			result[str] = append(result[str], str, lowV)
			continue
		}

		temp[ordered] = lowV
	}

	return result
}

func sortChars(s string) string {
	chars := strings.Split(s, "")
	sort.Strings(chars)
	return strings.Join(chars, "")
}
