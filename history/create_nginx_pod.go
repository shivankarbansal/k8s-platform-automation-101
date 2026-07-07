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
	// 1. Connection setup for Minikube
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
	// PART A: IDEMPOTENT CREATION LOGIC FOR NAMESPACE
	// =================================================================
	targetNamespace := "personal-platform-sandbox"

	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), targetNamespace, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("Platform Status: Namespace '%s' already exists.\n", targetNamespace)
	} else if errors.IsNotFound(err) {
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
		fmt.Printf("Unexpected API error checking namespace: %v\n", err)
		return
	}

	// =================================================================
	// PART B: EXPLICIT REMOVAL LOGIC FOR OLD NAMESPACE
	// =================================================================
	oldNamespace := "platform-sandbox"
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), oldNamespace, metav1.DeleteOptions{})
	if err == nil {
		fmt.Printf("Success! Explicitly requested deletion of old namespace: %s\n", oldNamespace)
	} else if errors.IsNotFound(err) {
		// Pass silently if it's already gone
	} else {
		fmt.Printf("Failed to delete old namespace: %v\n", err)
	}

	// =================================================================
	// PART C: IDEMPOTENT DEPLOYMENT LOGIC FOR NGINX POD
	// =================================================================
	podName := "nginx-web-server"
	fmt.Printf("\nPlatform Automation: Checking if Pod '%s' exists in '%s'...\n", podName, targetNamespace)

	// 1. CHECK: Try to get the pod inside our specific namespace
	_, err = clientset.CoreV1().Pods(targetNamespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if err == nil {
		// Pod already exists!
		fmt.Printf("Platform Status: Pod '%s' is already deployed. Skipping creation.\n", podName)
	} else if errors.IsNotFound(err) {
		// 2. CREATE: Define the Pod specification programmatically (Equivalent to Pod YAML)
		fmt.Printf("Platform Automation: Pod '%s' not found. Deploying it now...\n", podName)

		podSpec := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: podName,
				Labels: map[string]string{
					"app": "nginx-platform-demo",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx-container",
						Image: "nginx:latest",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
							},
						},
					},
				},
			},
		}

		// 3. EXECUTE: Send the POST request to create the pod in the specific targetNamespace
		_, err = clientset.CoreV1().Pods(targetNamespace).Create(context.TODO(), podSpec, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to deploy Pod: %v\n", err)
			return
		}

		fmt.Printf("Success! Pod '%s' has been created via Go.\n", podName)
	} else {
		fmt.Printf("Unexpected API error checking Pod: %v\n", err)
		return
	}

	fmt.Println("\nExecution finished. Platform is reconciled.")
}
