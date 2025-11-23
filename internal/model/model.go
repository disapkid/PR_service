package model

type User struct {
	Id string `db:"id"`
	Name string `db:"name"`
	IsActive bool `db:"is_active"`
	TeamId int64 `db:"team_id"`
}

func (u User) Values() []interface{} {
    return []interface{}{
        u.Id,
        u.Name,
        u.IsActive,
        u.TeamId,
    }
}