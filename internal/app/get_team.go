package app

import (
	"context"
	"fmt"
	desc "pr_service/pkg/gen"
)

func (s *Service) TeamGetGet(ctx context.Context, params desc.TeamGetGetParams) (desc.TeamGetGetRes, error) {
	if len(params.TeamName) == 0 {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprint("team Name is empty"),
			},
		}, fmt.Errorf("team name is empty")
	}

	teamName := params.TeamName

	team, name, err  := s.userRepository.GetTeam(ctx, teamName)
	if err != nil {
		return &desc.ErrorResponse{
			Error: desc.ErrorResponseError{
				Message: fmt.Sprintf("s.userRepository.GetTeam: %v", err),
			},
		}, fmt.Errorf("s.userRepository.GetTeam: %v", err)
	}

	res := desc.Team{}
	for _, val := range team {
		var temp desc.TeamMember
		temp.UserID = val.Id
		temp.Username = val.Name
		temp.IsActive = val.IsActive
		res.Members = append(res.Members, temp)
	}
	res.TeamName = name

	return &res, nil
}