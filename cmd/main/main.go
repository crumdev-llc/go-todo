package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	todo "github.com/crumdev/go-todo/api"
)

var todoFileName = ".todo.json"

func main() {

	add := flag.Bool("add", false, "Task to be included in the ToDo list")
	del := flag.Int("delete", 0, "Delete item N from Todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	//v := flag.Bool("v", false, "verbose output for list")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool, Developed from The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Define an items list
	l := &todo.List{}

	//Use the Get method to read to do items from file
	if err := l.GetFile(todoFileName); err != nil {
		fmt.Println(todoFileName + " not found...creating now.")
		os.Create(todoFileName)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	//for no extra arguyments , print the list
	case *list:
		fmt.Print(l)
	case *complete > 0:
		//complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		//save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		//add the task
		t, err := getTaskInput(os.Stdin, flag.Args()...)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.Add(t)

		//save the list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *del > 0:
		task, err := l.GetTask(*del)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Deleting item %d: %s\n", *del, task)
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		//save the list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid Option")
		os.Exit(1)
	}

}

func getTaskInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}
