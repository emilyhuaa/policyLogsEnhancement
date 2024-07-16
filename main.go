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

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			pods, err := pkg.GetPods(clientset)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			services, err := pkg.GetServices(clientset)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			cacheMutex.Lock()
			pkg.CacheMetadata(metadataCache, pods.Items, services.Items)
			cacheMutex.Unlock()
		}
	}()

	// Print the cache every 30 seconds
	go func() {
		for range ticker.C {
			cacheMutex.Lock()
			fmt.Println("IP Address : Name/Namespace")
			for ip, metadataList := range metadataCache {
				for _, metadata := range metadataList {
					fmt.Printf("%s : %s/%s\n", ip, metadata.Name, metadata.Namespace)
				}
			}
			cacheMutex.Unlock()
		}
	}()

	// Keep the main goroutine running
	select {}

}
