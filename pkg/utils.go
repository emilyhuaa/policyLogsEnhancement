package pkg

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Metadata struct {
	Name      string
	Namespace string
}

// GetPods lists Kubernetes Pods by namespace(s)
func GetPods(client kubernetes.Interface) (*v1.PodList, error) {
	pods, err := client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v", err)
		return nil, err
	}
	return pods, nil
}

// GetServices lists all Kubernetes Services across all namespaces
func GetServices(client kubernetes.Interface) (*v1.ServiceList, error) {
	services, err := client.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting services: %v", err)
		return nil, err
	}
	return services, nil
}

// CacheMetadata caches pod and service information by IP address
func CacheMetadata(metadataCache map[string][]Metadata, pods []v1.Pod, services []v1.Service) {
	// Create a set of IP addresses from the current pod and service data
	currentIPs := make(map[string]struct{})

	for _, pod := range pods {
		podIP := pod.Status.PodIP
		podName := pod.Name
		podNamespace := pod.Namespace

		if podIP == pod.Status.HostIP {
			continue
		}

		metadataCache[podIP] = []Metadata{{Name: podName, Namespace: podNamespace}}
		currentIPs[podIP] = struct{}{}
	}

	// Cache service metadata
	for _, srv := range services {
		srvIP := srv.Spec.ClusterIP
		srvName := srv.Name
		srvNamespace := srv.Namespace

		metadataCache[srvIP] = []Metadata{{Name: srvName, Namespace: srvNamespace}}
		currentIPs[srvIP] = struct{}{}
	}

	// Remove entries from the cache for IPs that are no longer present
	for ip := range metadataCache {
		if _, ok := currentIPs[ip]; !ok {
			delete(metadataCache, ip)
		}
	}
}
