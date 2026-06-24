package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
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
	UseSudo      bool   `env:"use_sudo"`
}

// Config holds validated step configuration.
type Config struct {
	file         string
	oldValue     string
	newValue     string
	showFile     bool
	notfoundExit bool
	useSudo      bool
}

// Step implements the change-value step logic.
type Step struct {
	inputParser stepconf.InputParser
	cmdFactory  command.Factory
	logger      log.Logger
}

// NewStep creates a Step with injected dependencies.
func NewStep(inputParser stepconf.InputParser, cmdFactory command.Factory, logger log.Logger) Step {
	return Step{
		inputParser: inputParser,
		cmdFactory:  cmdFactory,
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
		useSudo:      input.UseSudo,
	}, nil
}

// Run performs the value replacement in the target file.
func (s Step) Run(cfg Config) error {
	content, err := s.readFile(cfg.file, cfg.useSudo)
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

	if err := s.writeFile(cfg.file, []byte(replaced), cfg.useSudo); err != nil {
		return fmt.Errorf("failed to write file (%s): %w", cfg.file, err)
	}

	s.logger.Donef("Done")
	return nil
}

// readFile reads the target file, falling back to an elevated read (when
// useSudo is set) if the step lacks permission. On the 2026 Linux stack the
// step runs as the non-root `ubuntu` user, so root-owned files require sudo.
func (s Step) readFile(path string, useSudo bool) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err == nil {
		return content, nil
	}
	if !useSudo || !errors.Is(err, fs.ErrPermission) {
		return nil, err
	}

	s.logger.Warnf("Permission denied reading %s, retrying with sudo", path)
	var stdout, stderr bytes.Buffer
	cmd := s.cmdFactory.Create("sudo", []string{"-n", "cat", path}, &command.Opts{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if runErr := cmd.Run(); runErr != nil {
		return nil, fmt.Errorf("elevated read failed (%s): %w", strings.TrimSpace(stderr.String()), runErr)
	}
	return stdout.Bytes(), nil
}

// writeFile writes content to the target file, falling back to an elevated
// write (when useSudo is set) if the step lacks permission (see readFile for
// context). `sudo tee` preserves the existing file's owner and mode; `-n`
// prevents sudo from blocking on a password prompt (and from consuming the
// piped content).
func (s Step) writeFile(path string, content []byte, useSudo bool) error {
	err := os.WriteFile(path, content, 0644)
	if err == nil {
		return nil
	}
	if !useSudo || !errors.Is(err, fs.ErrPermission) {
		return err
	}

	s.logger.Warnf("Permission denied writing %s, retrying with sudo", path)
	var stderr bytes.Buffer
	cmd := s.cmdFactory.Create("sudo", []string{"-n", "tee", path}, &command.Opts{
		Stdin:  bytes.NewReader(content),
		Stdout: io.Discard,
		Stderr: &stderr,
	})
	if runErr := cmd.Run(); runErr != nil {
		return fmt.Errorf("elevated write failed (%s): %w", strings.TrimSpace(stderr.String()), runErr)
	}
	return nil
}

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)
	cmdFactory := command.NewFactory(envRepository)
	step := NewStep(inputParser, cmdFactory, logger)

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
