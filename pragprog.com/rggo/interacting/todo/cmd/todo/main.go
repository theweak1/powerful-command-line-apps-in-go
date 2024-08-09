package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"pragprog.com/rggo/interacting/todo"
)

// Hardcoding the file name
var todoFilename = ".todo.json"

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for the Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "\nInstructions:")
		fmt.Fprintln(flag.CommandLine.Output(), "  - To add a new task, use the '-add' flag followed by the task description.")
		fmt.Fprintln(flag.CommandLine.Output(), "    Example: ./todo -add \"Buy groceries\"")
		fmt.Fprintln(flag.CommandLine.Output(), "    Example: ./todo -add walk the dog")
		fmt.Fprintln(flag.CommandLine.Output(), "  - To list all tasks, use the '-list' flag.")
		fmt.Fprintln(flag.CommandLine.Output(), "  - To list all pending tasks, use the '-p' flag.")
		fmt.Fprintln(flag.CommandLine.Output(), "  - To mark a task as complete, use the '-complete' flag followed by the task number.")
		fmt.Fprintln(flag.CommandLine.Output(), "    Example: ./todo -complete 1")
		fmt.Fprintln(flag.CommandLine.Output(), "  - To delete a task, use the '-del' flag followed by the task number.")
		fmt.Fprintln(flag.CommandLine.Output(), "    Example: ./todo -del 2")
		fmt.Fprintln(flag.CommandLine.Output(), "  - To view tasks with additional details (such as creation date), use the '-v' flag.")
		fmt.Fprintln(flag.CommandLine.Output(), "\nEnvironment Variables:")
		fmt.Fprintln(flag.CommandLine.Output(), "  - You can set the TODO_FILENAME environment variable to specify a custom file name for the todo list.")
	}

	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("del", 0, "Item to be deleted from list")
	verbose := flag.Bool("v", false, "verbose view")
	pending := flag.Bool("p", false, "List pending tasks")

	flag.Parse()

	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFilename = os.Getenv("TODO_FILENAME")
	}

	// Define an items list
	l := &todo.List{}

	// Use the Get method to read to do items from file
	if err := l.Get(todoFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what todo based on the number of arguments provided
	switch {
	// For no extra arguments, print the list
	case *list:
		// List current todo items
		fmt.Print(l)

	case *verbose:
		// List current todo items with created date/time
		fmt.Print(l.Verbose())

	case *pending:
		// List all pending todo items
		fmt.Print(l.Pend())

	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		// When any arguments (excluding flags) are provided they will be
		// used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Add the task
		for _, task := range strings.Split(t, "\n") {
			l.Add(task)
		}

		// Save the new list
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// getTask function decides where to get the description for a new
// task from: arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	var tasks []string
	s := bufio.NewScanner(r)

	// Scan each line of input
	for s.Scan() {
		line := s.Text()
		if len(line) > 0 {
			tasks = append(tasks, line)
		}
	}

	// Check for scanning errors
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	// Join all tasks into a single string separated by newline characters
	return strings.Join(tasks, "\n"), nil
}
