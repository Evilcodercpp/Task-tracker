package transportgrpc

import (
	"net"

	taskpb "github.com/Evilcodercpp/project-protos/proto/task"
	"github.com/Evilcodercpp/task-service/internal/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPC(svc *task.Service) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}
	grpcSrv := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcSrv, NewHandler(svc))
	reflection.Register(grpcSrv)
	return grpcSrv.Serve(lis)
}
