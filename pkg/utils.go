package pkg

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodInfo struct {
	Name      string
	Namespace string
}

// ListPods lists Kubernetes Pods by namespace(s)
func ListPods(namespace string, client kubernetes.Interface) (*v1.PodList, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v", err)
		return nil, err
	}
	return pods, nil
}

// CachePods caches pod information by node IP address
func CachePods(podInfoCache map[string][]PodInfo, pods []v1.Pod) {
	// Create a set of IP addresses from the current cache
	currentIPs := make(map[string]bool)
	for ip := range podInfoCache {
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

		podInfo := PodInfo{
			Name:      podName,
			Namespace: podNamespace,
		}

		// Set the pod information for the current IP
		podInfoCache[podIP] = []PodInfo{podInfo}

		// Remove the IP from the set of current IPs
		delete(currentIPs, podIP)
	}

	// Remove any remaining IPs from the cache
	for ip := range currentIPs {
		delete(podInfoCache, ip)
	}
}
