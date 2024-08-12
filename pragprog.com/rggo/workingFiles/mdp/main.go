package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"

	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
<p><strong>File:</strong> {{ .FileName }}</p>
{{ .Body }}
  </body>
</html>
`
)

type content struct {
	Title    string
	Body     template.HTML
	FileName string
}

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	if err := run(*filename, *tFname, os.Stdin, os.Stdout, *skipPreview); err != nil {

		os.Exit(1)
	}
}

func run(filename, tFname string, r io.Reader, out io.Writer, skipPreview bool) error {
	// Check if an environment variable is set for the template
	if tFname == "" {
		tFname = os.Getenv("MDP_TEMPLATE")
	}

	// Read all the data from the input file and check for errors
	var input []byte
	var err error

	// If a filename is provided, read from the file
	if filename != "" {
		input, err = os.ReadFile(filename)
		if err != nil {
			return err
		}
	} else {
		// Otherwise, read from STDIN
		input, err = io.ReadAll(r)
		if err != nil {
			return err
		}
	}

	// Parse the content
	htmlData, err := parseContent(input, tFname, filename)
	if err != nil {
		return err
	}

	// Create temporary file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}

	// Close the temp file after creation
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)

	// Save the generated HTML to the temporary file
	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	// Ensure to delete the `Outname` file after the run function has completed
	defer os.Remove(outName)

	// Preview the generated HTML
	return preview(outName)
}

func parseContent(input []byte, tFname, fileName string) ([]byte, error) {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// If user provided alternate template file, replace template
	if tFname != "" {
		t, err = template.ParseFiles(tFname)

		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("File is: %s", fileName)
	// Instantiate the content type, adding the title and body
	c := content{
		Title:    "Markdown Preview Tool",
		Body:     template.HTML(body),
		FileName: fileName,
	}

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Execute the template with the content type
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outFname, data, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)

	if err != nil {
		return err
	}

	// Open the file using default program
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)
	return err
}
