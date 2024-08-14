package main

import (
	"flag"
	"fmt"
	"io"

	"log"
	"os"
	"path/filepath"
)

type config struct {
	// Extension to filter out
	ext string
	// Minimum file size for processing
	size int64
	// Flag to list files only
	list bool
	// Flag to delete files
	del bool
	// Log destination writer
	wLog io.Writer
	// Directory to archive files to
	archive string
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start")   // Root directory for the file operations
	logFile := flag.String("log", "", "Log deletes to this file") // Optional log file to record deletions
	// Action options
	list := flag.Bool("list", false, "List files only")        // Option to list files without taking further actions
	del := flag.Bool("del", false, "Delete files")             // Option to delete files
	archive := flag.String("archive", "", "Archive directory") // Directory where files should be archived
	// Filter options
	ext := flag.String("ext", "", "File extension to filter out") // File extension to filter (exclude)
	size := flag.Int64("size", 0, "Minimum file size")            // Minimum file size for processing
	flag.Parse()                                                  // Parse the command line flags

	var (
		f   = os.Stdout // Default to stdout for logging
		err error       // Error handling variable
	)

	// If a log file is specified, open it for appending logs
	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1) // Exit if log file cannot be opened
		}
		defer f.Close() // Ensure the log file is closed when the function returns
	}

	// Configuring the application based on parsed flags
	c := config{
		ext:     *ext,
		size:    *size,
		list:    *list,
		del:     *del,
		wLog:    f,
		archive: *archive,
	}

	// Run the main file processing function
	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1) // Exit if there is an error during execution
	}
}

// run processes files based on the provided configuration
func run(root string, out io.Writer, cfg config) error {
	// Initialize the logger for deleted files
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)

	// Walk through the directory tree starting from the root
	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err // Return error if one occurs while walking the file tree
			}

			// Filter out files based on extension and size criteria
			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil // Skip the file if it doesn't meet the criteria
			}

			// If list option is enabled, just list the file and return
			if cfg.list {
				return listFile(path, out)
			}

			// If an archive directory is specified, archive the file
			if cfg.archive != "" {
				if err := ArchiveFile(cfg.archive, root, path); err != nil {
					return err // Return error if archiving fails
				}
			}

			// If delete option is enabled, delete the file
			if cfg.del {
				return delFile(path, delLogger)
			}

			// Default action is to list the file if no other actions were taken
			return listFile(path, out)
		})
}

