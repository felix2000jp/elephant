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

Elephant uses a **feature-based architecture** with a flat structure that leverages Bubble Tea's native event system. Each feature contains components that handle different states and share common messages.

```
internal/
├── core/                        # Domain logic (pure Go, no UI dependencies)
│   ├── note.go                  # Note entity
│   ├── repository.go            # File system operations
│   └── repository_test.go       # Repository tests
├── features/                    # Feature modules
│   └── notes/                   # Notes feature
│       ├── feature.go           # Feature orchestrator with state management
│       ├── list_component.go    # List notes component
│       ├── view_component.go    # View note component
│       ├── edit_component.go    # Edit note component
├── theme/                       # UI styling
│   └── style.go                 # Application styles and themes
└── app/                         # Application orchestrator
    └── model.go                 # Main model & message router
```

### Architecture Principles

#### Component Structure

- Each component (list, view, edit) is a separate file with its own responsibilities
- The `feature.go` file orchestrates state transitions and component coordination
- Components have both `ForegroundUpdate` (when active) and `BackgroundUpdate` (when inactive) methods
- State management is centralized in the feature orchestrator

#### Message Flow

- App layer routes messages to the notes feature
- Feature orchestrator manages state transitions based on messages (ViewNoteMsg, EditNoteMsg, etc.)
- Components handle their own specific UI updates and user interactions
- Shared messages are defined in `messages.go` for cross-component communication

## Notes

- Standard Go conventions and best practices should be applied as development progresses
- Use the feature-based architecture for all new functionality
- Keep core domain logic in `internal/core/` with no UI dependencies
- Features should communicate through the app layer, not directly with each other