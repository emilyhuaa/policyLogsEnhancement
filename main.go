package main

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"

	pb "github.com/emilyhuaa/policyLogsEnhancement/pkg/rpc"
	utils "github.com/emilyhuaa/policyLogsEnhancement/pkg/utils"
	"github.com/go-logr/stdr"
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
	stdLogger := log.New(os.Stderr, "", log.LstdFlags)
	utils.Logger = stdr.New(stdLogger)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		utils.Logger.Error(err, "Error getting user home dir")
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		utils.Logger.Error(err, "Error getting kubernetes config")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		utils.Logger.Error(err, "Failed to create Kubernetes clientset")
		os.Exit(1)
	}

	go utils.UpdateCache(clientset, utils.Logger)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		utils.Logger.Error(err, "failed to listen")
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCacheServiceServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		utils.Logger.Error(err, "failed to serve")
		os.Exit(1)
	}
}
