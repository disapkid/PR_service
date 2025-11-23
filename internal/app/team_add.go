package app

import (
	"context"
	"fmt"

	"pr_service/internal/model"
	desc "pr_service/pkg/gen"
)

func (s *Service) TeamAddPost(ctx context.Context, req *desc.Team) (desc.TeamAddPostRes, error) {
	if req == nil {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprint("req is nil"),
			},
		}, fmt.Errorf("req is nil")
	}


	modelTeam := model.Team{
		Name: req.TeamName,
	}

	modelUsers := make([]*model.User, 0, len(req.Members))
	for _, user := range req.Members {
		modelUsers = append(modelUsers, &model.User{
			Name: user.Username,
			Id: user.UserID,
			IsActive: user.IsActive,
		})
	}

	_, err := s.userRepository.CreateTeam(ctx, modelTeam, modelUsers)
	if err != nil {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprintf("s.userRepository.CreateTeam: %v", err),
			},
		}, fmt.Errorf("s.userRepository.CreateTeam: %v", err)
	}

	return &desc.TeamAddPostCreated{
		Team: desc.NewOptTeam(desc.Team{
			TeamName: req.TeamName,
			Members: req.Members,
		}),
	}, nil
}
