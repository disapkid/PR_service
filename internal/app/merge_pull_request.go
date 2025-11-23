package app

import (
	"context"
	"fmt"
	desc "pr_service/pkg/gen"
)

func (s *Service) PullRequestMergePost(ctx context.Context, req *desc.PullRequestMergePostReq) (desc.PullRequestMergePostRes, error) {
	if req == nil {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprint("req is nil"),
			},
		}, fmt.Errorf("req is nil")
	}

	prId := req.PullRequestID

	err := s.userRepository.MergePullRequest(ctx, prId)
	if err != nil {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprintf("s.userRepository.MergePullRequest: %v", err),
			},
		}, fmt.Errorf("s.userRepository.MergePullRequest: %v", err)
	}

	return &desc.PullRequestMergePostOK{
		Pr: desc.NewOptPullRequest(desc.PullRequest{}),
	}, nil
}