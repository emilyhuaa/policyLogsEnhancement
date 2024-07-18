package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	utils "github.com/emilyhuaa/policyLogsEnhancement/pkg"
	pb "github.com/emilyhuaa/policyLogsEnhancement/pkg/rpc"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type server struct {
	pb.UnimplementedCacheServiceServer
}

func (s *server) GetCache(ctx context.Context, req *pb.CacheRequest) (*pb.CacheResponse, error) {
	cache := utils.GetMetadataCache()

	var cacheData []*pb.IPMetadata
	for ip, metadata := range cache {
		cacheData = append(cacheData, &pb.IPMetadata{
			Ip: ip,
			Metadata: &pb.Metadata{
				Name:      metadata.Name,
				Namespace: metadata.Namespace,
			},
		})
	}
	return &pb.CacheResponse{Data: cacheData}, nil
}

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

	go utils.UpdateCache(clientset)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCacheServiceServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
