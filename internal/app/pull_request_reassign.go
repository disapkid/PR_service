package app

import (
	"context"
	"fmt"

	"pr_service/internal/model"
	desc "pr_service/pkg/gen"
)

func (s *Service) PullRequestReassignPost(ctx context.Context, req *desc.PullRequestReassignPostReq) (desc.PullRequestReassignPostRes, error)  {
	if req == nil {
		return &desc.PullRequestReassignPostConflict{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprint("req is nil"),
			},
		}, fmt.Errorf("req is nil")
	}


	modelRev := model.Reviewer{
		PullRequestId: req.GetPullRequestID(),
		UserId: req.GetOldUserID(),
	}

	err := s.userRepository.PullRequestReassign(ctx, modelRev)
	if err != nil {
		// need fix
		return &desc.PullRequestReassignPostConflict{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprintf("s.userRepository.PullRequestReassign: %v", err),
			},
		}, fmt.Errorf("s.userRepository.PullRequestReassign: %v", err)
	}

	return &desc.PullRequestReassignPostOK{
		// need fix
	}, nil
}