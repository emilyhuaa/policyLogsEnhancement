package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/emilyhuaa/policyLogsEnhancement/pkg"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	metadataCache := make(map[string][]pkg.Metadata)
	var cacheMutex sync.Mutex

	// Start a goroutine to update the cache every 10 seconds
	podTicker := time.NewTicker(10 * time.Second)
	go func() {
		for range podTicker.C {
			pods, err := pkg.GetPods(clientset)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			cacheMutex.Lock()
			pkg.CachePods(metadataCache, pods.Items)
			cacheMutex.Unlock()
		}
	}()

	srvTicker := time.NewTicker(30 * time.Second)
	go func() {
		for range srvTicker.C {
			services, err := pkg.GetServices(clientset)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			cacheMutex.Lock()
			pkg.CacheServices(metadataCache, services.Items)
			cacheMutex.Unlock()
		}
	}()

	// Print the cache every 10 seconds
	// go func() {
	// 	for range ticker.C {
	// 		cacheMutex.Lock()
	// 		fmt.Println("Pod IP Address : Pod Name/Namespace")
	// 		for ip, podInfoList := range metadataCache {
	// 			for _, podInfo := range podInfoList {
	// 				fmt.Printf("%s : %s/%s\n", ip, podInfo.Name, podInfo.Namespace)
	// 			}
	// 		}
	// 		cacheMutex.Unlock()
	// 	}
	// }()

	// Keep the main goroutine running
	select {}

}
