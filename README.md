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
bashsim$ kubectl get deployments
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   1/1     1            1           12s
bashsim$ kubectl get pods
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-574b87c764-x9z2p   1/1     Running   0          25s
bashsim$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                      READY   STATUS    RESTARTS   AGE
default       nginx-deployment-574b87c764-x9z2p         1/1     Running   0          45s
kube-system   coredns-787d4945fb-k2j4l                  1/1     Running   0          2d
kube-system   etcd-minikube                             1/1     Running   0          2d
kube-system   kube-apiserver-minikube                   1/1     Running   0          2d
kube-system   kube-controller-manager-minikube          1/1     Running   0          2d
kube-system   kube-proxy-5j8r2                          1/1     Running   0          2d
kube-system   kube-scheduler-minikube                   1/1     Running   0          2d
kube-system   storage-provisioner                       1/1     Running   0          2d
bashsim$ kubectl get ns
NAME              STATUS   AGE
default           Active   2d
kube-node-lease   Active   2d
kube-public       Active   2d
kube-system       Active   2d
bashsim$
bashsim$
bashsim$
bashsim$ kubectl create ns ab
namespace/ab created
bashsim$ kubectl get ns ab -o yaml
apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: "2023-10-27T10:01:30Z"
  name: ab
  resourceVersion: "492"
  uid: 84a3c1b2-d4e5-4f6a-9b8c-7d0e1f2a3b4c
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
bashsim$
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
