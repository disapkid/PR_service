package app

import (
	"context"
	"fmt"
	desc "pr_service/pkg/gen"
)


func (s *Service) UsersSetIsActivePost(ctx context.Context, req *desc.UsersSetIsActivePostReq) (desc.UsersSetIsActivePostRes, error) {
	if req == nil {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprint("req is nil"),
			},
		}, fmt.Errorf("req is nil")
	}

	usrId := req.UserID
	status := req.IsActive
	
	err := s.userRepository.SetUserFlag(ctx, usrId, status)
	if err != nil {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprintf("s.userRepository.SetUserFlag: %v", err),
			},
		}, fmt.Errorf("s.userRepository.SetUserFlag: %v", err)
	}

	return &desc.UsersSetIsActivePostOK{
		User: desc.NewOptUser(desc.User{}),
	}, nil
}