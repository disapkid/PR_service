package model

type PullRequest struct {
	Id string `db:"id"`
	Name string `db:"name"`
	AuthorId string `db:"author_id"`
	Status string `db:"status"`
}

func (pr PullRequest) Values() []interface{} {
	return []interface{} {
		pr.Id,
		pr.Name,
		pr.AuthorId,
		pr.Status,
	}
}