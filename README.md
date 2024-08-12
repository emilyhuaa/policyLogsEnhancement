# aws-kubernetes-metadata-cache-deployment
Amazon EKS Metadata Cache deployment maintains a cache with current pod and service metadata information. This information can be used to enhance logs, debug, and troubleshoot.

## Getting Started
You’ll need a Kubernetes cluster version 1.25+ to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.

## Setup
Download the latest version of the [yaml](https://github.com/emilyhuaa/policyLogsEnhancement/blob/main/aws-k8s-metadata.yaml) and apply it to the cluster.
```kubectl apply -f aws-k8s-metadata.yaml```

Alternatively, downloading the eks/aws-vpc/cni helm chart will automatically download all necessary components.

