package main

import (
	"flag"
	"fmt"
	"os"

	todo "github.com/crumdev/go-todo/api"
)

var todoFileName = ".todo.json"

func main() {

	task := flag.String("task", "", "Task to be included in the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
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
	if err := l.Get(todoFileName); err != nil {
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

	case *task != "":
		//add the task
		l.Add(*task)

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
