package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Session struct {
	FilePath string
	History  []Entry
}

func New(path string) (*Session, error) {
	if path == "" {
		return nil, fmt.Errorf("session path cannot be empty")
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	s := &Session{
		FilePath: path,
		History:  []Entry{},
	}

	// Start reading if file exists
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open session file: %w", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	for decoder.More() {
		var entry Entry
		if err := decoder.Decode(&entry); err != nil {
			// If we encounter an error, we stop reading but return what we have so far
			// or we could return error. Given it's a history file, maybe warning is better?
			// But for now let's be strict to avoid corruption issues being ignored.
			return nil, fmt.Errorf("failed to decode entry: %w", err)
		}
		s.History = append(s.History, entry)
	}

	return s, nil
}

func (s *Session) Append(input, output string) error {
	if strings.TrimSpace(input) == "" && strings.TrimSpace(output) == "" {
		return nil
	}

	entry := Entry{
		Input:  input,
		Output: output,
	}

	s.History = append(s.History, entry)

	f, err := os.OpenFile(s.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open session file for append: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(entry); err != nil {
		return fmt.Errorf("failed to encode entry: %w", err)
	}

	return nil
}
