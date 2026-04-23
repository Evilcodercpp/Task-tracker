package transportgrpc

import (
	"context"
	"fmt"

	taskpb "github.com/Evilcodercpp/project-protos/proto/task"
	userpb "github.com/Evilcodercpp/project-protos/proto/user"
	"github.com/Evilcodercpp/task-service/internal/task"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	svc        *task.Service
	userClient userpb.UserServiceClient
	taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *task.Service, uc userpb.UserServiceClient) *Handler {
	return &Handler{svc: svc, userClient: uc}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	t := &task.Task{Title: req.Title, UserID: uint(req.UserId)}
	if err := h.svc.CreateTask(t); err != nil {
		return nil, status.Errorf(codes.Internal, "create task: %v", err)
	}
	return &taskpb.CreateTaskResponse{Task: toProto(t)}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	t, err := h.svc.GetTaskByID(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get task: %v", err)
	}
	return &taskpb.GetTaskResponse{Task: toProto(t)}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, total, err := h.svc.GetAllTasks(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list tasks: %v", err)
	}
	pbTasks := make([]*taskpb.Task, len(tasks))
	for i := range tasks {
		pbTasks[i] = toProto(&tasks[i])
	}
	return &taskpb.ListTasksResponse{Tasks: pbTasks, Total: uint32(total)}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksByUserResponse, error) {
	tasks, _, err := h.svc.GetTasksByUserID(uint(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list tasks by user: %v", err)
	}
	pbTasks := make([]*taskpb.Task, len(tasks))
	for i := range tasks {
		pbTasks[i] = toProto(&tasks[i])
	}
	return &taskpb.ListTasksByUserResponse{Tasks: pbTasks}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	updated, err := h.svc.UpdateTaskByID(uint(req.Id), task.Task{Title: req.Title, IsDone: req.IsDone})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update task: %v", err)
	}
	return &taskpb.UpdateTaskResponse{Task: toProto(updated)}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if err := h.svc.DeleteTaskByID(uint(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "delete task: %v", err)
	}
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func toProto(t *task.Task) *taskpb.Task {
	return &taskpb.Task{
		Id:     uint32(t.ID),
		Title:  t.Title,
		UserId: uint32(t.UserID),
		IsDone: t.IsDone,
	}
}
