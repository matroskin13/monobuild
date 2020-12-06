package service

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

type Writer struct {
	packageName string
}

func NewWriter(packageName string) *Writer {
	return &Writer{packageName: packageName}
}

func (w Writer) Write(p []byte) (n int, err error) {
	fmt.Println(fmt.Sprintf("[%s]", strings.ToUpper(w.packageName)), string(p))

	return len(p), nil
}

type Service struct {
	entry  string
	pid    int
	cmd    *exec.Cmd
	logger io.Writer
	stop   chan struct{}
}

func NewService(entry string, logger io.Writer) (*Service, error) {
	s := &Service{entry: entry, logger: logger, stop: make(chan struct{})}

	if err := s.start(entry); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) OnStop() chan struct{} {
	return s.stop
}

func (s *Service) Stop() error {
	return syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
}

func (s *Service) Reload() error {
	if err := s.Stop(); err != nil {
		return fmt.Errorf("cannot stop service: %w", err)
	}

	return s.start(s.entry)
}

func (s *Service) start(entry string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Dir = entry
	cmd.Stdout = s.logger
	cmd.Stderr = s.logger

	if err := cmd.Start(); err != nil {
		return err
	}

	s.cmd = cmd
	s.pid = cmd.Process.Pid

	return nil
}
