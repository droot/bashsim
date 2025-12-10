# bashsim

A Bash shell simulator powered by Google's Gemini LLM.

## Description

`bashsim` simulates a bash shell environment using the Gemini Large Language Model. It maintains context across commands by storing session history, allowing for a realistic and persistent shell experience. It supports standard shell features like command history navigation (readline), heredoc syntax, and positional parameters.

## Install

### Prerequisites
- Go 1.24 or later
- A valid Gemini API Key

### Build from Source

```bash
# Clone the repository (if applicable)
# git clone https://github.com/droot/bashsim.git
# cd bashsim

# Build the binary
go build -o bashsim .
```

## Usage

Before running, ensure you have your Gemini API Key exported:

```bash
export GEMINI_API_KEY="your-api-key-here"
# OR
export GOOG_API_KEY="your-api-key-here"
```

### Interactive Mode (REPL)

Start the simulator:
```bash
./bashsim
```

**Features:**
- **Context Awareness**: The LLM remembers previous commands in the session.
- **History Navigation**: Use **Up/Down** arrow keys to cycle through previous commands.
- **History Search**: Press **Ctrl+R** to search your command history.
- **Heredoc Support**: You can use multi-line input using `<<EOF` syntax.

### One-Shot Execution

Run a single command string and exit:
```bash
./bashsim -c "echo 'Hello from bashsim'"
```

You can also pass positional parameters (simulating `$0`, `$1`, etc.):
```bash
# Simulates: $0=script_name, $1=foo
./bashsim -c "echo \$0 argument is \$1" script_name foo
```

### Configuration

You can configure `bashsim` using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `BASHSIM_SESSION` | Path to the session history file (JSONL). | `/tmp/bashsim.session.default` |
| `BASHSIM_MODEL` | Gemini model name to use. | `gemini-3-pro-preview` |
