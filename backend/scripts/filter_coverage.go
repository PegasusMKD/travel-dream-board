// scripts/filter_coverage.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run filter_coverage.go <file> <pattern1> [pattern2] [pattern3]...")
		os.Exit(1)
	}

	filePath := os.Args[1]
	excludePatterns := os.Args[2:]

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		shouldExclude := false

		for _, pattern := range excludePatterns {
			if strings.Contains(line, pattern) {
				shouldExclude = true
				break
			}
		}

		if !shouldExclude {
			lines = append(lines, line)
		}
	}

	os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
