package model

type Team struct {
	Id int64 `db:"id"`
	Name string `db:"name"`
}