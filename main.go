package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/droot/bashsim/llm"
	"github.com/droot/bashsim/session"

	"github.com/chzyer/readline"
)

func main() {
	// 0. Parse Flags
	// We use NewFlagSet to avoid polluting global state if we add tests later,
	// but mostly to control parsing behavior.
	fs := flag.NewFlagSet("bashsim", flag.ExitOnError)
	cmdFlag := fs.String("c", "", "command to execute")
	fs.Parse(os.Args[1:])

	// 1. Session Path
	sessionPath := os.Getenv("BASHSIM_SESSION")
	if sessionPath == "" {
		sessionPath = "/tmp/bashsim.session.default"
	}

	// 2. Initialize Session
	sess, err := session.New(sessionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing session: %v\n", err)
		os.Exit(1)
	}

	// 3. Initialize LLM
	ctx := context.Background()
	// Using a default model, e.g., gemini-3-pro-preview.
	// We could also make this configurable via env var if needed.
	modelName := os.Getenv("BASHSIM_MODEL")
	if modelName == "" {
		modelName = "gemini-3-pro-preview"
	}

	client, err := llm.New(ctx, modelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing LLM client: %v\nPlease ensure GEMINI_API_KEY is set.\n", err)
		os.Exit(1)
	}

	// 4. Handle -c option
	if *cmdFlag != "" {
		handleCommandMode(ctx, client, sess, *cmdFlag, fs.Args())
		return
	}

	// 5. REPL
	runREPL(ctx, client, sess)
}

func handleCommandMode(ctx context.Context, client *llm.Client, sess *session.Session, cmd string, args []string) {
	// If arguments are present, we should inform the LLM about them.
	// bash -c 'echo $0' arg0
	// args[0] is $0
	input := cmd
	if len(args) > 0 {
		var sb strings.Builder
		sb.WriteString("Context: The following positional parameters are set:\n")
		for i, arg := range args {
			sb.WriteString(fmt.Sprintf("$%d=%s\n", i, arg))
		}
		sb.WriteString("Command to execute:\n")
		sb.WriteString(cmd)
		input = sb.String()
	}

	resp, err := client.GenerateResponse(ctx, sess.History, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating response: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(resp)
	if len(resp) > 0 && resp[len(resp)-1] != '\n' {
		fmt.Println()
	}

	// We probably don't want to save "Context: ..." to history in the exact verbose form,
	// or maybe we do to keep context. For now, saving as is.
	if err := sess.Append(cmd, resp); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
	}
}

func runREPL(ctx context.Context, client *llm.Client, sess *session.Session) {
	historyFile := filepath.Join(os.Getenv("HOME"), ".bashsim_history")
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "bashsim$ ",
		HistoryFile:       historyFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	var heredocDelim string
	var fullInput strings.Builder

	for {
		// Update prompt based on state
		if len(heredocDelim) > 0 {
			rl.SetPrompt("> ")
		} else {
			rl.SetPrompt("bashsim$ ")
		}

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		if len(heredocDelim) > 0 {
			fullInput.WriteString(line + "\n")
			if strings.TrimSpace(line) == heredocDelim {
				// Heredoc finished
				input := fullInput.String()
				heredocDelim = ""
				fullInput.Reset()
				processInput(ctx, client, sess, input)
			}
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.TrimSpace(line) == "exit" {
			break
		}

		// Check for heredoc start
		if delim, ok := parseHeredoc(line); ok {
			heredocDelim = delim
			fullInput.WriteString(line + "\n")
			continue
		}

		processInput(ctx, client, sess, line)
	}
}

func processInput(ctx context.Context, client *llm.Client, sess *session.Session, input string) {
	// Generate response
	resp, err := client.GenerateResponse(ctx, sess.History, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating response: %v\n", err)
		return
	}

	// Output response
	fmt.Print(resp)
	// Ensure newline at end if not present to mimic shell output cleanly if missing
	if len(resp) > 0 && resp[len(resp)-1] != '\n' {
		fmt.Println()
	}

	// Save to session
	if err := sess.Append(input, resp); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
	}
}

func parseHeredoc(line string) (string, bool) {
	// Simple parsing for << DELIM or << "DELIM" or << 'DELIM'
	// This is not a full bash parser but sufficient for simulation
	idx := strings.Index(line, "<<")
	if idx == -1 {
		return "", false
	}

	rest := strings.TrimSpace(line[idx+2:])
	if rest == "" {
		return "", false
	}

	// Handle quotes
	if strings.HasPrefix(rest, "\"") && strings.Contains(rest[1:], "\"") {
		end := strings.Index(rest[1:], "\"")
		return rest[1 : end+1], true
	}
	if strings.HasPrefix(rest, "'") && strings.Contains(rest[1:], "'") {
		end := strings.Index(rest[1:], "'")
		return rest[1 : end+1], true
	}

	// Simple word
	parts := strings.Fields(rest)
	if len(parts) > 0 {
		return parts[0], true
	}

	return "", false
}
