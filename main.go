package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. Automatically locate and build kubeconfig across Windows, Mac, or Linux
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfigLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfigLoader.ClientConfig()
	if err != nil {
		fmt.Printf("Failed to load kubeconfig from default paths: %v\n", err)
		return
	}

	// 2. Create the Kubernetes Clientset interface
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

	_, err = clientset.CoreV1().Pods(targetNamespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if err == nil {
		fmt.Printf("Platform Status: Pod '%s' is already deployed. Skipping creation.\n", podName)
	} else if errors.IsNotFound(err) {
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

	// =================================================================
	// PART D: INSPECTING POD HEALTH, STATUS AND LIVE IP ADDRESS
	// =================================================================
	fmt.Printf("\nPlatform Inspection: Fetching real-time status of '%s'...\n", podName)

	livePod, err := clientset.CoreV1().Pods(targetNamespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Failed to fetch live pod details: %v\n", err)
		return
	}

	podPhase := livePod.Status.Phase
	podIP := livePod.Status.PodIP
	if podIP == "" {
		podIP = "Not Assigned Yet"
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("Pod Metadata Name : %s\n", livePod.Name)
	fmt.Printf("Current Run Phase : %s\n", podPhase)
	fmt.Printf("Assigned Pod IP   : %s\n", podIP)
	fmt.Println("----------------------------------------")

	// =================================================================
	// PART E: CLUSTER-WIDE INVENTORY DISCOVERY (FOR LOOP)
	// =================================================================
	fmt.Println("\nPlatform Discovery: Scanning ALL pods across ALL namespaces...")

	podList, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to fetch global cluster pod inventory: %v\n", err)
		return
	}

	fmt.Printf("Found %d active system/application pods running in total.\n", len(podList.Items))
	fmt.Println("=======================================================================")
	fmt.Printf("%-30s %-30s %-15s\n", "NAMESPACE", "POD NAME", "STATUS")
	fmt.Println("=======================================================================")

	for _, pod := range podList.Items {
		fmt.Printf("%-30s %-30s %-15s\n", pod.Namespace, pod.Name, pod.Status.Phase)
	}
	fmt.Println("=======================================================================")

	fmt.Println("\nExecution finished. Platform is reconciled.")
}
