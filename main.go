package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/emilyhuaa/policyLogsEnhancement/pkg/rpc"
	utils "github.com/emilyhuaa/policyLogsEnhancement/pkg/utils"
	"github.com/go-logr/stdr"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	config, err := rest.InClusterConfig()
	if err != nil {
		utils.Logger.Error(err, "Failed to get in-cluster config")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.Logger.Error(err, "Failed to create Kubernetes clientset")
		os.Exit(1)
	}

	go utils.UpdateCache(clientset)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		utils.Logger.Error(err, "failed to listen")
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCacheServiceServer(grpcServer, &server{})
	utils.Logger.Info("server listening", "address", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		utils.Logger.Error(err, "failed to serve")
		os.Exit(1)
	}
}
