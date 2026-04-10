package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: qn-new \"Title of note\"")
		os.Exit(1)
	}

	var title strings.Builder
	for i := 1; i < len(os.Args); i++ {
		if i > 1 {
			title.WriteString(" ")
		}
		title.WriteString(os.Args[i])
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dir := os.Getenv("QN_DIR")
	if dir == "" {
		dir = filepath.Join(home, "quicknotes")
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	dtg := now.Format("20060102150405")
	filename := dtg + ".md"
	fullpath := filepath.Join(dir, filename)

	file, err := os.Create(fullpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content := fmt.Sprintf("# %s\n\nCreated: %s\nTags:\n\n---\n\n", title, now.Format(time.RFC3339))

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("pwsh.exe", "-Command", editor, fullpath)
	} else {
		cmd = exec.Command(editor, fullpath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
