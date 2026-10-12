# ctracker-CLI

![Go version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)

## Demo

<!-- Replace this placeholder with your demo GIF. -->
![ctracker CLI demo](YOUR_DEMO_GIF_HERE.gif)

## About

ctracker is a small command-line tool that helps you see the Git projects in a folder and browse their commit history.

## Features

- Find Git repositories inside a folder and its subfolders.
- Show each project name, commit count, and local path.
- Open a project to browse its commit messages.
- Move through the project list and commit history in a terminal interface.

## Tech Stack

| Area | Technology |
| --- | --- |
| Language | Go 1.26 |
| Terminal UI | Bubble Tea, Bubbles, Lip Gloss |
| Version control | Git |
| Testing | No automated tests currently |
| Deployment | Not applicable; runs locally in your terminal |
| CI/CD | No CI/CD workflow is currently set up |

## Getting Started

### Requirements

- Go 1.26 or later
- Git installed and available in your terminal

### Setup

Clone the repository and build the program:

```bash
git clone https://github.com/juanroy012/ctracker-CLI.git
cd ctracker-CLI
go build -o ctracker .
```

## Usage

Run the program and enter the folder you want to scan:

```bash
./ctracker
```

Or pass the folder path as an argument:

```bash
./ctracker ~/Projects
```

Use the arrow keys to move through the project list. Press `enter` to view a project's commits, `esc` or `backspace` to return to the list, and `q` or `ctrl+c` to quit.

You can also run it without building a binary:

```bash
go run . ~/Projects
```

## Project Structure

- `main.go` starts the terminal interface and displays projects and commits.
- `project.go` finds Git repositories and reads their commit history.
- `util.go` handles terminal input and expands paths that start with `~`.

## Links

- Repository: [github.com/juanroy012/ctracker-CLI](https://github.com/juanroy012/ctracker-CLI)
- Contact: [GitHub profile](https://github.com/juanroy012)
