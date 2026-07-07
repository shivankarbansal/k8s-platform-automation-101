package history

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. Locate and build connection configuration for local Minikube
	homeDir, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(homeDir, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create clientset: %v\n", err)
		return
	}

	// =================================================================
	// PART A: IDEMPOTENT CREATION LOGIC FOR PERSONAL-PLATFORM-SANDBOX
	// =================================================================
	targetNamespace := "personal-platform-sandbox"

	// Check if the target namespace already exists
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), targetNamespace, metav1.GetOptions{})
	
	if err == nil {
		fmt.Printf("Platform Status: Namespace '%s' already exists. Skipping creation safely.\n", targetNamespace)
	} else if errors.IsNotFound(err) {
		// If it's truly not found, create it dynamically
		fmt.Printf("Platform Automation: '%s' not found. Creating it now...\n", targetNamespace)
		nsSpec := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: targetNamespace},
		}
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create namespace: %v\n", err)
			return
		}
		fmt.Println("Success! Namespace initialized gracefully.")
	} else {
		// Handle unexpected system errors (e.g., API server down)
		fmt.Printf("Unexpected API error checking namespace: %v\n", err)
		return
	}

	// =================================================================
	// PART B: EXPLICIT REMOVAL LOGIC FOR OLD PLATFORM-SANDBOX
	// =================================================================
	oldNamespace := "platform-sandbox"

	fmt.Printf("Platform Automation: Checking if old namespace '%s' needs removal...\n", oldNamespace)

	// Send a transactional DELETE request to the Minikube API [1]
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), oldNamespace, metav1.DeleteOptions{})
	
	if err == nil {
		fmt.Printf("Success! Explicitly requested deletion of old namespace: %s\n", oldNamespace)
	} else if errors.IsNotFound(err) {
		// If it was already deleted previously, pass silently without breaking the script
		fmt.Printf("Platform Status: Old namespace '%s' does not exist anymore. No deletion needed.\n", oldNamespace)
	} else {
		fmt.Printf("Failed to delete old namespace: %v\n", err)
	}

	fmt.Println("\nExecution finished. Platform is reconciled.")
}
