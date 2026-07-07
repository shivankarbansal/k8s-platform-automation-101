# Kubernetes Platform Automation Tool

A high-performance, cross-platform CLI automation utility written in Go (Golang) designed to interact dynamically with Kubernetes clusters via the official `client-go` SDK [INDEX]. This project shifts from traditional declarative YAML configurations to an imperative, transactional approach for managing infrastructure.

## 🚀 Key Features

*   **Platform Agnostic:** Automatically locates and builds cluster credentials (`kubeconfig`) across Windows, macOS, and Linux [INDEX] without hardcoded path dependencies.
*   **Idempotent Resource Management:** Verification checkpoints evaluate cluster states before execution, preventing duplicate resource errors or execution crashes.
*   **Automated Workload Injection:** Programmatically defines and provisions target Namespaces and containerized applications (Nginx Web Server Pods).
*   **Cluster Discovery Engine:** Scans global cluster environments to output an aligned, real-time inventory table of all active system and tenant pods.

---

## 🛠️ Architecture & Workflow

1.  **Authentication:** Evaluates system configuration variables to initialize a secure `Clientset` interface [INDEX].
2.  **Reconcile Namespace:** Evaluates if `personal-platform-sandbox` exists. If missing, spins it up dynamically.
3.  **Explicit Cleanup:** Targets and executes a transactional `DELETE` loop on legacy resources (`platform-sandbox`).
4.  **Telemetry & Inspection:** Deploys an application container workload and retrieves assigned IP network blocks.
5.  **Global Inventory Scan:** Queries the internal cluster database (`etcd`) and prints an active inventory report.

---

## 📋 Prerequisites & Setup

Ensure you have a local cluster running (e.g., [Minikube](https://k8s.io)) and **Go** installed on your workstation.

### 1. Initialize Module & Track Dependencies

```bash
# Initialize a new Go module scope
go mod init k8s-platform-automation

# Fetch the official Kubernetes core API dependencies
go get k8s.io/client-go@v0.35.1
go get k8s.io/apimachinery@v0.35.1
```

> **Observation:** Initializing the module generates the `go.mod` file to handle package definitions. Injecting core API dependencies generates the `go.sum` signature validation file for secure module hash tracking.

---

## 💻 Execution & Packaging Workflow

### 1. Reconcile Dependencies & Run Locally
In case you face dependency mismatches or missing module checksum flags when executing the script, clear and sync the tracking references:

```bash
# Clear old unused sub-modules and sync the new cross-platform imports
go mod tidy

# Force a clean recompile and execute the automation logic
go run -a main.go
```

### 2. Compile and Package into Standalone Binary
To skip runtime compilation lag, compile your script into a self-contained, native system binary executable with zero external runtime dependencies:

#### 🪟 Windows Setup (PowerShell)
```powershell
# Compile into a standalone Windows executable
go build -o platform-tool.exe main.go

# Execute the packaged engine binary instantly
.\platform-tool.exe
```

#### 🍏 macOS / 🐧 Linux Setup (Bash/Zsh)
```bash
# Compile into a native Unix binary
go build -o platform-tool main.go

# Grant execution permissions to the binary
chmod +x platform-tool

# Execute the packaged engine binary instantly
./platform-tool
```

---

## 🛡️ Production Code Structure

```text
k8s-platform-automation/
├── .gitignore          # Strict ignore tracking rules for binaries/local caches
├── go.mod              # Package dependencies definition index
├── go.sum              # SHA-256 cryptographic verification file
├── main.go             # Central Go platform engine logic source file
└── README.md           # Professional system operational blueprint
```

---

## 🧹 Systematic Environment Cleanup (Local Best Practices)

To mirror high-cost public cloud resource discipline on your local machine, run the teardown loop before powering down your workstation:

```bash
# 1. Nuke workloads and dynamic namespace blocks
kubectl delete namespace personal-platform-sandbox

# 2. Stop the local cluster virtualization engine
minikube stop

# 3. Verify clean system slate
minikube status
```
