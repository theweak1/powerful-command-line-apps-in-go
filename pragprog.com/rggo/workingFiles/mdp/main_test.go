package main

import (
	"bytes"
	"os"

	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := parseContent(input, "", inputFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer

	// Run the 'run' function with the test markdown file
	if err := run(inputFile, "", os.Stdin, &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	// Get the name of the output HTML file from the buffer
	resultFile := strings.TrimSpace(mockStdOut.String())

	// Read the generated HTML file
	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	// Read the expected HTML output (golden file)
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the generated HTML to the expected HTML
	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}

	// Clean up: remove the temporary file created by the run function
	os.Remove(resultFile)
}
