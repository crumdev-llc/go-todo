package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName      = "todo"
	todoFileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")
	windowsCheck()
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	if os.Getenv("TODO_FILENAME") != "" {
		path, err := filepath.EvalSymlinks(os.Getenv("TODO_FILENAME"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		todoFileName = path
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(todoFileName)
	os.Exit(result)
}

func formatList(l []string) string {
	formatted := ""
	prefix := "  "
	for k := range l {
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, l[k])
	}
	return formatted
}

func formatCompletedList(l []string) string {
	formatted := ""
	prefix := "X "
	for k := range l {
		if k == 0 {
			formatted += fmt.Sprintf("%s%d: %s\n", "  ", k+1, l[k])
		} else {
			formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, l[k])
		}
	}
	return formatted
}

func TestTodoCLI(t *testing.T) {
	var tasks [3]string
	tasks[0] = "test task number 1"
	tasks[1] = "test task number 2"
	tasks[2] = "test task number 3"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)
	t.Run("AddNewTask", func(t *testing.T) {
		for i := range tasks {
			cmd := exec.Command(cmdPath, "-task", tasks[i])
			if err := cmd.Run(); err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("ListTasks", func(t *testing.T) {

		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := formatList(tasks[0:])

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("CompleteTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "2")
		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		cmd1 := exec.Command(cmdPath, "-complete", "3")
		_, err1 := cmd1.CombinedOutput()
		if err1 != nil {
			t.Fatal(err1)
		}

		cmd2 := exec.Command(cmdPath, "-list")
		out, err2 := cmd2.CombinedOutput()
		if err2 != nil {
			t.Fatal(err2)
		}

		expected := formatCompletedList(tasks[0:])

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}

func windowsCheck() {
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
}
