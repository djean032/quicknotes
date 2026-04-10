package config

import (
	"os"
	"path/filepath"
)

func QNDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dir := os.Getenv("QN_DIR")

	if dir == "" {
		dir = filepath.Join(home, "quicknotes")
	}

	return dir
}

func Editor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}
	return editor
}

func IndexPath() string {
	dir := QNDir()
	dir = filepath.Join(dir, ".qn_index.json")
	return dir
}
