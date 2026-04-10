package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type project struct {
	name   string
	commit []string
	path   string
}

// InitializeProject take project path and return a project object
func initializeProject(folderPath string) (project, error) {
	output, err := exec.Command("git", "-C", folderPath, "log", "--oneline").Output()
	if err != nil {
		return project{}, err
	}
	commit := separateCommits(output)
	name := extractProjectName(folderPath)

	return project{
		name:   name,
		commit: commit,
		path:   folderPath,
	}, nil
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

func separateCommits(output []byte) []string {
	c := strings.Split(string(output), "\n")
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
			folders = append(folders, strings.TrimSuffix(path, "/.git"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return folders, nil
}
