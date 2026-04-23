package transportgrpc

import (
	"context"

	taskpb "github.com/Evilcodercpp/project-protos/proto/task"
	"github.com/Evilcodercpp/task-service/internal/task"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	svc *task.Service
	taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *task.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	t := &task.Task{Title: req.Title}
	if err := h.svc.CreateTask(t); err != nil {
		return nil, status.Errorf(codes.Internal, "create task: %v", err)
	}
	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{Id: uint32(t.ID), Title: t.Title},
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	t, err := h.svc.GetTaskByID(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get task: %v", err)
	}
	return &taskpb.GetTaskResponse{
		Task: &taskpb.Task{Id: uint32(t.ID), Title: t.Title},
	}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, total, err := h.svc.GetAllTasks(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list tasks: %v", err)
	}
	pbTasks := make([]*taskpb.Task, len(tasks))
	for i, t := range tasks {
		pbTasks[i] = &taskpb.Task{Id: uint32(t.ID), Title: t.Title}
	}
	return &taskpb.ListTasksResponse{
		Tasks: pbTasks,
		Total: uint32(total),
	}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	updated, err := h.svc.UpdateTaskByID(uint(req.Id), task.Task{Title: req.Title})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update task: %v", err)
	}
	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{Id: uint32(updated.ID), Title: updated.Title},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if err := h.svc.DeleteTaskByID(uint(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "delete task: %v", err)
	}
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}
