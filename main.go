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

	namespace := "" // An empty string returns all namespaces
	pods, err := pkg.ListPods(namespace, clientset)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	podInfoCache := pkg.CachePods(pods.Items)

	fmt.Println("Pod IP Address : Pod Name/Namespace:")
	for ip, podInfoList := range podInfoCache {
		for _, podInfo := range podInfoList {
			fmt.Printf("%s : %s,%s\n", ip, podInfo.Name, podInfo.Namespace)
		}
	}

	ipAddress := "192.168.77.216"

	if podInfoList, ok := podInfoCache[ipAddress]; ok {
		fmt.Printf("Pod IP Address: %s\n", ipAddress)
		for _, podInfo := range podInfoList {
			fmt.Printf("Pod Name: %s, Namespace: %s\n", podInfo.Name, podInfo.Namespace)
		}
	} else {
		fmt.Printf("No pods found for IP address: %s\n", ipAddress)
	}

}
