package utils

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Metadata struct {
	Name      string
	Namespace string
}

var (
	MetadataCache = make(map[string]Metadata)
	Logger        logr.Logger
)
var CacheMutex sync.Mutex

// CacheMetadata caches pod and service information by IP address
func UpdateMetadataCache(metadataCache map[string]Metadata, pods []v1.Pod, services []v1.Service) {
	// Create a set of IP addresses from the current pod and service data
	currentIPs := make(map[string]struct{})

	for _, pod := range pods {
		podIP := pod.Status.PodIP

		if podIP == pod.Status.HostIP {
			metadataCache[podIP] = Metadata{Name: "hostIP", Namespace: "hostIP"}
			currentIPs[podIP] = struct{}{}
		} else {
			metadataCache[podIP] = Metadata{Name: pod.Name, Namespace: pod.Namespace}
			currentIPs[podIP] = struct{}{}
		}
	}

	// Cache service metadata
	for _, srv := range services {
		srvIP := srv.Spec.ClusterIP

		metadataCache[srvIP] = Metadata{Name: srv.Name, Namespace: srv.Namespace}
		currentIPs[srvIP] = struct{}{}
	}

	// Remove entries from the cache for IPs that are no longer present
	for ip := range metadataCache {
		if _, ok := currentIPs[ip]; !ok {
			delete(metadataCache, ip)
		}
	}
}

// // CacheMetadata caches pod and service information by IP address
// func UpdateMetadataCache(metadataCache map[string]Metadata, pods []v1.Pod, services []v1.Service) {
//     // Create a slice to store current IP addresses
//     var currentIPs []string

//     for _, pod := range pods {
//         podIP := pod.Status.PodIP

//         if podIP == pod.Status.HostIP {
//             metadataCache[podIP] = Metadata{Name: "hostIP", Namespace: "hostIP"}
//         } else {
//             metadataCache[podIP] = Metadata{Name: pod.Name, Namespace: pod.Namespace}
//         }

//         // Add the IP address to the currentIPs slice if it doesn't exist
//         if !contains(currentIPs, podIP) {
//             currentIPs = append(currentIPs, podIP)
//         }
//     }

//     // Cache service metadata
//     for _, srv := range services {
//         srvIP := srv.Spec.ClusterIP

//         metadataCache[srvIP] = Metadata{Name: srv.Name, Namespace: srv.Namespace}

//         // Add the IP address to the currentIPs slice if it doesn't exist
//         if !contains(currentIPs, srvIP) {
//             currentIPs = append(currentIPs, srvIP)
//         }
//     }

//     // Remove entries from the cache for IPs that are no longer present
//     for ip := range metadataCache {
//         if !contains(currentIPs, ip) {
//             delete(metadataCache, ip)
//         }
//     }
// }

// // contains checks if a string is present in a slice
// func contains(slice []string, str string) bool {
//     for _, s := range slice {
//         if s == str {
//             return true
//         }
//     }
//     return false
// }

func UpdateCache(client kubernetes.Interface) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				Logger.Error(err, "Error getting pods")
				continue
			}
			services, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				Logger.Error(err, "Error getting services")
				continue
			}
			CacheMutex.Lock()
			UpdateMetadataCache(MetadataCache, pods.Items, services.Items)
			CacheMutex.Unlock()
			Logger.Info("Updated Metadata Cache", "Cache", MetadataCache)
		}
	}
}

func GetMetadataCache() map[string]Metadata {
	CacheMutex.Lock()
	defer CacheMutex.Unlock()
	return MetadataCache
}
