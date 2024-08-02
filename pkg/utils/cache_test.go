package utils

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCacheMetadata(t *testing.T) {
	// Create test data
	pods := []v1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod1",
				Namespace: "default",
			},
			Status: v1.PodStatus{
				PodIP:  "10.0.0.1",
				HostIP: "10.0.0.2",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod2",
				Namespace: "kube-system",
			},
			Status: v1.PodStatus{
				PodIP:  "10.0.0.3",
				HostIP: "10.0.0.3",
			},
		},
	}

	services := []v1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service1",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "10.0.0.4",
			},
		},
	}

	expectedCache := map[string]Metadata{
		"10.0.0.1": {Name: "pod1", Namespace: "default"},
		"10.0.0.3": {Name: "hostIP", Namespace: "hostIP"},
		"10.0.0.4": {Name: "service1", Namespace: "default"},
	}

	// Call the function under test
	UpdateMetadataCache(MetadataCache, pods, services)

	// Check if the cache is populated correctly
	if !reflect.DeepEqual(MetadataCache, expectedCache) {
		t.Errorf("CacheMetadata() failed. Expected: %v, Got: %v", expectedCache, MetadataCache)
	}

	// Test case: Pod deleted
	podsAfterDeletion := []v1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod2",
				Namespace: "kube-system",
			},
			Status: v1.PodStatus{
				PodIP:  "10.0.0.3",
				HostIP: "10.0.0.3",
			},
		},
	}

	expectedCacheAfterDeletion := map[string]Metadata{
		"10.0.0.3": {Name: "hostIP", Namespace: "hostIP"},
		"10.0.0.4": {Name: "service1", Namespace: "default"},
	}

	UpdateMetadataCache(MetadataCache, podsAfterDeletion, services)

	// Check if the cache is updated correctly after pod deletion
	if !reflect.DeepEqual(MetadataCache, expectedCacheAfterDeletion) {
		t.Errorf("CacheMetadata() failed after pod deletion. Expected: %v, Got: %v", expectedCacheAfterDeletion, MetadataCache)
	}
}
