package model

type Event struct {
	ID          int64  `db:"id"`
	Owner       int64  `db:"owner"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
