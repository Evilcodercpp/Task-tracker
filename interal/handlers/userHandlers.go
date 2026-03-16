package handlers

import (
	"context"
	"errors"

	userservice "Task-tracker/interal/userService"
	"Task-tracker/interal/web/users"

	"github.com/google/uuid"
)

type UserHandler struct {
	svc *userservice.UserService
}

func NewUserHandler(svc *userservice.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUsers(_ context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	all, err := h.svc.GetAllUsers()
	if err != nil {
		return nil, err
	}

	response := users.GetUsers200JSONResponse{}
	for _, u := range all {
		id := u.ID
		email := u.Email
		password := u.Password
		response = append(response, users.User{
			Id:       &id,
			Email:    &email,
			Password: &password,
		})
	}
	return response, nil
}

func (h *UserHandler) PostUsers(_ context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	if request.Body.Email == nil {
		return nil, errors.New("field 'email' is required")
	}
	if request.Body.Password == nil {
		return nil, errors.New("field 'password' is required")
	}

	user := &userservice.User{
		ID:       uuid.NewString(),
		Email:    *request.Body.Email,
		Password: *request.Body.Password,
	}

	if err := h.svc.CreateUser(user); err != nil {
		return nil, err
	}

	id := user.ID
	email := user.Email
	password := user.Password

	return users.PostUsers201JSONResponse{
		Id:       &id,
		Email:    &email,
		Password: &password,
	}, nil
}

func (h *UserHandler) GetUsersId(_ context.Context, request users.GetUsersIdRequestObject) (users.GetUsersIdResponseObject, error) {
	user, err := h.svc.GetUserByID(request.Id)
	if err != nil {
		return nil, err
	}

	id := user.ID
	email := user.Email
	password := user.Password

	return users.GetUsersId200JSONResponse{
		Id:       &id,
		Email:    &email,
		Password: &password,
	}, nil
}

func (h *UserHandler) PatchUsersId(_ context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersIdResponseObject, error) {
	existing, err := h.svc.GetUserByID(request.Id)
	if err != nil {
		return nil, err
	}

	if request.Body.Email != nil {
		existing.Email = *request.Body.Email
	}
	if request.Body.Password != nil {
		existing.Password = *request.Body.Password
	}

	if err := h.svc.UpdateUser(existing); err != nil {
		return nil, err
	}

	id := existing.ID
	email := existing.Email
	password := existing.Password

	return users.PatchUsersId200JSONResponse{
		Id:       &id,
		Email:    &email,
		Password: &password,
	}, nil
}

func (h *UserHandler) DeleteUsersId(_ context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	if err := h.svc.DeleteUser(request.Id); err != nil {
		return nil, err
	}
	return users.DeleteUsersId204Response{}, nil
}
