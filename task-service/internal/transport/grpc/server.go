package transportgrpc

import (
	"net"

	taskpb "github.com/Evilcodercpp/project-protos/proto/task"
	userpb "github.com/Evilcodercpp/project-protos/proto/user"
	"github.com/Evilcodercpp/task-service/internal/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPC(svc *task.Service, uc userpb.UserServiceClient) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}
	grpcSrv := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcSrv, NewHandler(svc, uc))
	reflection.Register(grpcSrv)
	return grpcSrv.Serve(lis)
}
