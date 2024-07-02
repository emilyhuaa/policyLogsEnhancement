package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/emilyhuaa/policyLogsEnhancement/pkg"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		fmt.Printf("error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}
	// An empty string returns all namespaces
	namespace := "kube-system"
	pods, err := pkg.ListPods(namespace, clientset)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, pod := range pods.Items {
		fmt.Printf("Pod name: %v\n", pod.Name)
	}
	var message string
	if namespace == "" {
		message = "Total Pods in all namespaces"
	} else {
		message = fmt.Sprintf("Total Pods in namespace `%s`", namespace)
	}
	fmt.Printf("%s %d\n", message, len(pods.Items))

	//ListNamespaces function call returns a list of namespaces in the kubernetes cluster
	namespaces, err := pkg.ListNamespaces(clientset)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
	}
	fmt.Printf("Total namespaces: %d\n", len(namespaces.Items))

	podIPAddresses, err := pkg.ListPodIPAddresses(namespace, clientset)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, ip := range podIPAddresses {
		fmt.Println(ip)
	}

}
