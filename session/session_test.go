package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSession(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "bashsim_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	sessionPath := filepath.Join(tmpDir, "test.session")

	// Test New with non-existent file
	s, err := New(sessionPath)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if len(s.History) != 0 {
		t.Errorf("Expected empty history, got %d", len(s.History))
	}

	// Test Append
	if err := s.Append("echo hello", "hello"); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Test New with existing file
	s2, err := New(sessionPath)
	if err != nil {
		t.Fatalf("New (existing) failed: %v", err)
	}
	if len(s2.History) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(s2.History))
	}
	if s2.History[0].Input != "echo hello" {
		t.Errorf("Expected input 'echo hello', got '%s'", s2.History[0].Input)
	}
}
