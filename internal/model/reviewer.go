package model

type Reviewer struct {
	PullRequestId string `db:"pr_id"`
	UserId string `db:"user_id"`
}

func (r Reviewer) Values() []interface{} {
	return []interface{} {
		r.PullRequestId,
		r.UserId,
	}
}