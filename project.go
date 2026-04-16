package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type project struct {
	name   string
	commit []string
	path   string
}

// InitializeProject take project path and return a project object
func initializeProject(folderPath string, ch chan project, wg *sync.WaitGroup) {
	defer wg.Done()
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("git", "-C", folderPath, "log", "--oneline")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := stderr.String()
		if strings.Contains(strings.ToLower(msg), "does not have any commits yet") {
			fmt.Printf("Directory %v doesn't have any commits yet\n", folderPath)
		} else {
			fmt.Printf("Skipping directory: %v likely not a repo or bad path\n", folderPath)
		}
	}

	name := extractProjectName(folderPath)
	commit := separateCommits(stdout.String())

	pr := project{
		name:   name,
		commit: commit,
		path:   folderPath,
	}
	ch <- pr
}

func extractProjectName(folderPath string) string {
	n := filepath.Base(folderPath)
	if n == ".git" || n == "." {
		wd, err := os.Getwd()
		if err != nil {
			return "unknown"
		}
		n = filepath.Base(wd)
	}
	return n
}

func separateCommits(output string) []string {
	c := strings.Split(output, "\n")
	return c
}

// scanGitDir return path collection that has .git
func scanGitDir(folderPath string) ([]string, error) {
	var folders []string
	folderPath = strings.TrimSuffix(folderPath, "/")
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == "volumes" {
			return filepath.SkipDir
		}
		if d.IsDir() && d.Name() == ".git" {
			basePath, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error: %v", err)
			}
			if folderPath == "." {
				folders = append(folders, basePath+strings.TrimSuffix(path, ".git"))
			} else {
				folders = append(folders, strings.TrimSuffix(path, ".git"))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return folders, nil
}
