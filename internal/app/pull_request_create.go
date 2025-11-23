package app

import (
	"context"
	"fmt"

	"pr_service/internal/model"
	desc "pr_service/pkg/gen"
)

func (s *Service) PullRequestCreatePost(ctx context.Context, req *desc.PullRequestCreatePostReq) (desc.PullRequestCreatePostRes, error)  {
	if req == nil {
		return &desc.PullRequestCreatePostConflict{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprint("req is nil"),
			},
		}, fmt.Errorf("req is nil")
	}


	modelPr := model.PullRequest{
		Id: req.GetPullRequestID(),
		Name: req.GetPullRequestName(),
		AuthorId: req.GetAuthorID(),
		Status: "OPEN",
	}

	err := s.userRepository.CreatePullRequest(ctx, modelPr)
	if err != nil {
		// need fix
		return &desc.PullRequestCreatePostNotFound{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprintf("s.userRepository.CreatePullRequest: %v", err),
			},
		}, fmt.Errorf("s.userRepository.CreatePullRequest: %v", err)
	}

	return &desc.PullRequestCreatePostCreated{
		Pr: desc.OptPullRequest{
			// need fix
		},
	}, nil
}
