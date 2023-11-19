package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please provide file for processing")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error at opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	maxMemory := 0

	for scanner.Scan() {
		line := scanner.Text()
		memoryMatch := extractMemory(line)
		if memoryMatch != nil {
			memoryValue, err := strconv.Atoi(memoryMatch[1])
			if err == nil && memoryValue > maxMemory {
				maxMemory = memoryValue
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error on reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("max memory found: %d (kB)\n", maxMemory)
}

func extractMemory(line string) []string {
	re := regexp.MustCompile(`memory - (\d+) \(kB\)`)
	return re.FindStringSubmatch(line)
}
