package repository

var (
	TeamsTable = "teams"
	UsersTable = "users"
	PullRequestTable = "pull_requests"
	ReviewersTable = "reviewers"

	TeamsColumns = []string{"name"}
	UsersColumns = []string{"id", "name", "is_active", "team_id"}
	PullRequestColumns = []string{"id","name", "author_id", "status"}
	ReviewerColumns = []string{"pr_id", "user_id"}
)