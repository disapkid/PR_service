package app

import (
	"context"
	"fmt"
	desc "pr_service/pkg/gen"
)

func (s *Service) UsersGetReviewGet(ctx context.Context, params desc.UsersGetReviewGetParams) (*desc.UsersGetReviewGetOK, error) {
	if len(params.UserID) == 0 {
		return &desc.UsersGetReviewGetOK{}, fmt.Errorf("team name is empty")
	}

	userId := params.UserID
	reviews, err := s.userRepository.GetReview(ctx, userId)
	if err != nil {
		return &desc.UsersGetReviewGetOK{}, fmt.Errorf("team name is empty")
	}

	var res desc.UsersGetReviewGetOK
	res.UserID = userId

	for _, val := range reviews {
		var temp desc.PullRequestShort
		temp.AuthorID = val.AuthorId
		temp.PullRequestID = val.Id
		temp.PullRequestName = val.Name
		temp.Status = desc.PullRequestShortStatus(val.Status)

		res.PullRequests = append(res.PullRequests, temp)
	}

	return &res, nil
}
