package main

import (
	"compress/gzip" // Package for reading and writing compressed files in gzip format

	"fmt"           // Package for formatted I/O operations
	"io"            // Package for basic I/O interfaces
	"log"           // Package for logging
	"os"            // Package for file system operations
	"path/filepath" // Package for manipulating filename paths
)

// filterOut is a helper function that filters files based on certain conditions.
// It returns true if the file should be filtered out (i.e., excluded), and false otherwise.
func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize { // Exclude directories and files smaller than minSize
		return true
	}

	if ext != "" && filepath.Ext(path) != ext { // Exclude files that don't match the specified extension
		return true
	}
	return false
}

// listFile prints the file path to the provided output writer (e.g., stdout or a file).
func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

// delFile deletes the file at the given path and logs the deletion using the provided logger.
func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil { // Attempt to remove the file
		return err
	}
	delLogger.Println(path) // Log the deletion
	return nil
}

// ArchiveFile compresses a file and saves it to a specified destination directory.
// The resulting file is stored in gzip format.
func ArchiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir) // Check the status of the destination directory
	if err != nil {
		return err
	}

	if !info.IsDir() { // Ensure the destination is a directory
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// Determine the relative directory path to maintain the folder structure
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	// Create the target file path with a .gz extension
	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relDir, dest)

	// Create the target directory structure if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Open the target file for writing
	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close() // Ensure the file is closed when the function returns

	// Open the source file for reading
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close() // Ensure the file is closed when the function returns

	// Create a new gzip writer
	zw := gzip.NewWriter(out)
	zw.Name = filepath.Base(path)             // Set the original filename in the gzip metadata
	if _, err = io.Copy(zw, in); err != nil { // Copy the content of the source file to the gzip writer
		return err
	}

	// Close the gzip writer to finalize the gzip file
	if err := zw.Close(); err != nil {
		return err
	}

	return out.Close() // Close the target file
}

