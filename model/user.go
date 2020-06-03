package model

type User struct {
	ID     int64  `db:"id"`
	Cookie string `db:"cookie"`
}
