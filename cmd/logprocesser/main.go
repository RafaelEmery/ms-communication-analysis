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
	sumMemory := 0
	numMemoryValues := 0
	var sumDatabaseTime float64
	numDatabaseValues := 0

	for scanner.Scan() {
		line := scanner.Text()
		memoryMatch := extractMemory(line)
		if memoryMatch != nil {
			memoryValue, err := strconv.Atoi(memoryMatch[1])
			if err == nil {
				sumMemory += memoryValue
				numMemoryValues++
			}
		}
		dbTimeMatch := extractDatabaseTime(line)
		if dbTimeMatch != nil {
			dbTimeValue, err := strconv.ParseFloat(dbTimeMatch[1], 64)
			if err == nil {
				sumDatabaseTime += dbTimeValue
				numDatabaseValues++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error on reading file: %v\n", err)
		os.Exit(1)
	}

	averageMemory := float64(sumMemory) / float64(numMemoryValues)
	fmt.Printf("1 - memory average: %.2f (kB)\n", averageMemory)

	averageDatabaseTime := sumDatabaseTime / float64(numDatabaseValues)
	fmt.Printf("2 - database interaction time average: %.2f (ms)\n", averageDatabaseTime)
}

func extractMemory(line string) []string {
	re := regexp.MustCompile(`memory - (\d+) \(kB\)`)
	return re.FindStringSubmatch(line)
}

func extractDatabaseTime(line string) []string {
	re := regexp.MustCompile(`database interaction time - (\d+\.\d+)ms`)
	return re.FindStringSubmatch(line)
}
