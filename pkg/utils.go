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

// CachePods caches pod information by node IP address
func CachePods(metadataCache map[string][]Metadata, pods []v1.Pod) {
	// Create a set of IP addresses from the current cache
	currentIPs := make(map[string]bool)
	for ip := range metadataCache {
		currentIPs[ip] = true
	}

	for _, pod := range pods {
		podIP := pod.Status.PodIP
		podName := pod.Name
		podNamespace := pod.Namespace
		nodeIP := pod.Status.HostIP

		if podIP == nodeIP {
			continue
		}

		podInfo := Metadata{
			Name:      podName,
			Namespace: podNamespace,
		}

		// Set the pod information for the current IP
		metadataCache[podIP] = []Metadata{podInfo}

		// Remove the IP from the set of current IPs
		delete(currentIPs, podIP)
	}

	// Remove any remaining IPs from the cache
	for ip := range currentIPs {
		delete(metadataCache, ip)
	}
}

// CacheServices caches pod information by node IP address
func CacheServices(metadataCache map[string][]Metadata, services []v1.Service) {
	// Create a set of IP addresses from the current cache
	currentIPs := make(map[string]bool)
	for ip := range metadataCache {
		currentIPs[ip] = true
	}

	for _, srv := range services {
		srvIP := srv.Spec.ClusterIP
		srvName := srv.Name
		srvNamespace := srv.Namespace

		srvInfo := Metadata{
			Name:      srvName,
			Namespace: srvNamespace,
		}

		// Set the pod information for the current IP
		metadataCache[srvIP] = []Metadata{srvInfo}

		// Remove the IP from the set of current IPs
		delete(currentIPs, srvIP)
	}

	// Remove any remaining IPs from the cache
	for ip := range currentIPs {
		delete(metadataCache, ip)
	}
}
