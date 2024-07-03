package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	podInfoCache := make(map[string][]string)
	for _, pod := range pods.Items {
		podIP := pod.Status.PodIP
		podName := pod.Name
		podNamespace := pod.Namespace
		nodeIP := pod.Status.HostIP

		if podIP == nodeIP {
			continue
		}

		podInfo := fmt.Sprintf("%s/%s", podName, podNamespace)
		podInfoCache[podIP] = append(podInfoCache[podIP], podInfo)
	}

	fmt.Println("Pod IP Address : Pod Name/Namespace:")
	for ip, info := range podInfoCache {
		fmt.Printf("%s : %s\n", ip, info)
	}

	ipAddress := "192.168.77.216"

	if podInfoList, ok := podInfoCache[ipAddress]; ok {
		fmt.Printf("Pod IP Address: %s\n", ipAddress)
		for _, podInfo := range podInfoList {
			parts := strings.Split(podInfo, "/")
			podN := parts[0]
			podNS := parts[1]
			fmt.Printf("Pod Name: %s, Namespace: %s\n", podN, podNS)
		}
	} else {
		fmt.Printf("No pods found for IP address: %s\n", ipAddress)
	}

}
