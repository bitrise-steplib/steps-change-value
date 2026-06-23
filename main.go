package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

// Input holds raw step inputs parsed from environment variables.
type Input struct {
	File         string `env:"file,required"`
	OldValue     string `env:"old_value,required"`
	NewValue     string `env:"new_value"`
	ShowFile     bool   `env:"show_file"`
	NotfoundExit bool   `env:"notfound_exit"`
}

// Config holds validated step configuration.
type Config struct {
	file         string
	oldValue     string
	newValue     string
	showFile     bool
	notfoundExit bool
}

// Step implements the change-value step logic.
type Step struct {
	inputParser stepconf.InputParser
	logger      log.Logger
}

// NewStep creates a Step with injected dependencies.
func NewStep(inputParser stepconf.InputParser, logger log.Logger) Step {
	return Step{
		inputParser: inputParser,
		logger:      logger,
	}
}

// ProcessConfig parses and validates step inputs.
func (s Step) ProcessConfig() (Config, error) {
	var input Input
	if err := s.inputParser.Parse(&input); err != nil {
		return Config{}, fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	return Config{
		file:         input.File,
		oldValue:     input.OldValue,
		newValue:     input.NewValue,
		showFile:     input.ShowFile,
		notfoundExit: input.NotfoundExit,
	}, nil
}

// Run performs the value replacement in the target file.
func (s Step) Run(cfg Config) error {
	content, err := os.ReadFile(cfg.file)
	if err != nil {
		return fmt.Errorf("failed to read file (%s): %w", cfg.file, err)
	}

	if cfg.showFile {
		s.logger.Printf("")
		s.logger.Printf("------------------------------------------")
		s.logger.Printf("-------------OLD  FILE--------------------")
		s.logger.Printf("------------------------------------------")
		s.logger.Printf("%s", string(content))
		s.logger.Printf("------------------------------------------")
	}

	if cfg.notfoundExit && !strings.Contains(string(content), cfg.oldValue) {
		return fmt.Errorf("old value (%s) was not found in the file (%s)", cfg.oldValue, cfg.file)
	}

	s.logger.Printf("Replacing...")
	replaced := strings.ReplaceAll(string(content), cfg.oldValue, cfg.newValue)

	if cfg.showFile {
		s.logger.Printf("")
		s.logger.Printf("------------------------------------------")
		s.logger.Printf("-------------NEW  FILE--------------------")
		s.logger.Printf("------------------------------------------")
		s.logger.Printf("%s", replaced)
		s.logger.Printf("------------------------------------------")
	}

	if err := os.WriteFile(cfg.file, []byte(replaced), 0644); err != nil {
		return fmt.Errorf("failed to write file (%s): %w", cfg.file, err)
	}

	s.logger.Donef("Done")
	return nil
}

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)
	step := NewStep(inputParser, logger)

	cfg, err := step.ProcessConfig()
	if err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	if err := step.Run(cfg); err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	return 0
}
