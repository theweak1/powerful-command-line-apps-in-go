package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Define a boolean flag -l to count lines instead of words
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")

	// Parsing the flags provided by the user
	flag.Parse()

	// Get the remaining command-Line arguments (files)
	files := flag.Args() // This captures any filenames provided after the flags

	//If no files are provided, us STDIN
	if len(files) == 0 {
		out, err := count("", os.Stdin, *lines, *bytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		fmt.Println(out)
	} else {

		total := 0

		// process each file
		for _, file := range files {
			out, err := count(file, nil, *lines, *bytes)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processig file %s: %s\n", file, err)
				continue
			}
			fmt.Printf("%s: %d\n", file, out)
			total += out
		}
		// Print the total count if more than one file is provided
		if len(files) > 1 {
			fmt.Printf("Total: %d\n", total)
		}
	}
}

func count(Fname string, r io.Reader, countLines, countBytes bool) (int, error) {
	var reader io.Reader

	reader = r

	// Check if a file is being provided
	if Fname != "" {
		data, err := os.ReadFile(Fname)
		if err != nil {
			return 0, err
		}
		reader = bytes.NewReader(data)
	}

	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(reader)

	// If the count lines flag is not set, we want to count words so we define the scanner split type to words (default is split by lines)
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	// If the count bytes flag is set, we want to count bytes so we define the scanner split type to bytes (default is split by lines)
	if countBytes {
		scanner.Split(bufio.ScanBytes)
	}

	// Define a counter
	wc := 0

	// For every word or line scanned, add 1 to the counter
	for scanner.Scan() {
		wc++
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// Return the total
	return wc, nil
}
