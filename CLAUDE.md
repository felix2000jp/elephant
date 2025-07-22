# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Elephant is a terminal user interface (TUI) for viewing and editing notes. Upon running the application, the app will
look for a ".elephant" directory where .md files will be found, these files are the user notes.

The app has three different states:

1. List state where a list of notes with title (the filename) and description (the first line of the file, usually in
   Markdown file format).
2. View state, whereupon using the correct shortcut, the selected file contents are shown in Markdown format.
3. Edit state where after using the correct shortcut, a new view is shown with the selected file content is rendered on
   a textarea and the user is allowed to edit them.

This application is built using the Go programming language and libraries: Bubble Tea, Lipgloss, Glamour, and Bubbles.

## Development Commands

```bash
go build # Build the project

go fmt ./... # Format code
go vet ./... # Vet code for issues

go test ./... # Run tests
go test -v ./... # Run tests with verbose output
go test -run TestName # Run a specific test

go mod tidy # Initialize and download dependencies
```

## Project Structure

Elephant uses a **converged feature-based architecture** that leverages Bubble Tea's native event system. Each feature
is self-contained with a three-part pattern: messages → handlers → model+component.

```
internal/
├── core/                        # Domain logic (pure Go, no UI dependencies)
│   ├── note.go                  # Note entity
│   └── repository.go            # File system operations
├── features/                    # Feature modules
│   ├── notes/                   
│   │   ├── list/                # List notes feature
│   │   │   ├── messages.go      # Messages
│   │   │   ├── handlers.go      # Message handlers
│   │   │   ├── component.go     # Feature model + component
│   │   │   └── list_test.go     # Feature tests
│   │   ├── view/                # View note feature
│   │   │   ├── messages.go      # Messages
│   │   │   ├── handlers.go      # Message handlers
│   │   │   ├── component.go     # Feature model + component
│   │   │   └── view_test.go     # Feature tests
│   │   ├── edit/                # Edit note feature
│   │   │   ├── messages.go      # Messages
│   │   │   ├── handlers.go      # Message handlers
│   │   │   ├── component.go     # Feature model + component
│   │   │   └── edit_test.go     # Feature tests
│   │   └── messages.go          # Shared note messages (Edit, View)
└── app/                         # Application orchestrator
    └── model.go                 # Main model & message router
```

### Architecture Principles

#### Message Flow

- App layer routes messages to appropriate features
- Features handle their own messages via handlers
- Cross-feature communication happens through the app layer
- Shared concerns (like window resize) are forwarded to all relevant features

## Notes

- Standard Go conventions and best practices should be applied as development progresses
- Use the feature-based architecture for all new functionality
- Keep core domain logic in `internal/core/` with no UI dependencies
- Features should communicate through the app layer, not directly with each other