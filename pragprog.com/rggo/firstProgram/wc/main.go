package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Define a boolean flag -l to count lines instead of words
	lines := flag.Bool("l", false, "Count lines")

	bytes := flag.Bool("b",false, "Count bytes")
	// Parsing the flags provided by the user
	flag.Parse()

	// Calling the count function to count the number of words (or lines)
	// received from the Standard Input and printing it out
	fmt.Println(count(os.Stdin, *lines, *bytes))
}

func count(r io.Reader, countLines, countBytes bool) int {
	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)

	// If the count lines flag is not set, we want to count words so we define the scanner split type to words (default is split by lines)
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	if countBytes {
		scanner.Split(bufio.ScanBytes)
	}

	// Define a counter
	wc := 0

	// For every word or line scanned, add 1 to the counter
	for scanner.Scan() {
		wc++
	}

	// Return the total
	return wc
}
