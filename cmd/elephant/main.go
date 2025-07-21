package main

import (
	"elephant/internal/app"
	"log"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	setupLogging()

	model := app.NewModel()
	program := tea.NewProgram(&model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		slog.Error("Something went wrong, exiting...", "err", err)
		os.Exit(1)
	}
}

func setupLogging() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
