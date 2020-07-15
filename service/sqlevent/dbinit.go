package sqlevent

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sir-farfan/hack-a-bot/model"
)

var DBDriver = "sqlite3"

type Storage struct {
	DB *sqlx.DB
}

func New() *Storage {
	db, err := sqlx.Open(DBDriver, "hackabot.db")
	if err != nil {
		log.Printf("ERROR: connecting to the database: %v\n", err)
		return nil
	}

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS event (id INTEGER PRIMARY KEY, owner INTEGER, name TEXT, description TEXT)")
	_, err = statement.Exec()
	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, cookie TEXT)")
	_, err = statement.Exec()
	if err != nil {
		log.Printf("error creating database: %v\n", err)
	}

	return &Storage{DB: db}
}

func (s *Storage) GetEvent(owner int64) []model.Event {
	events := []model.Event{}
	statement := "SELECT * FROM event"
	if owner > 0 {
		statement += " WHERE owner = ?"
	}

	err := s.DB.Select(&events, statement, owner)
	if err != nil {
		log.Println(err)
	}

	return events
}

func (s *Storage) GetEventByID(id int64) []model.Event {
	events := []model.Event{}
	statement := "SELECT * FROM event  WHERE id = ?"

	err := s.DB.Select(&events, statement, id)
	if err != nil {
		log.Println(err)
	}

	return events
}

func (s *Storage) CreateEvent(event model.Event) error {
	insert, err := s.DB.Prepare("INSERT INTO event (owner, name, description) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Error during insert event: %v\n", err)
		return err
	}
	_, err = insert.Exec(event.Owner, event.Name, event.Description)

	return err
}

func (s *Storage) UpdateEvent(event model.Event) error {
	statement := "UPDATE event SET "
	if event.Name != "" {
		statement += "name=? WHERE id=?"
		update, _ := s.DB.Prepare(statement)
		_, err := update.Exec(event.Name, event.ID)
		return err
	}
	if event.Description != "" {
		statement += "description=? WHERE id=?"
		update, _ := s.DB.Prepare(statement)
		_, err := update.Exec(event.Description, event.ID)
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(id int64) error {
	delete, err := s.DB.Prepare("DELETE FROM event WHERE id=?")
	if err != nil {
		return err
	}
	_, err = delete.Exec(id)
	return err
}

func (s *Storage) UserCookieCreate(user model.User) error {
	insert, err := s.DB.Prepare("INSERT INTO user (id, cookie) VALUES (?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = insert.Exec(user.ID, user.Cookie)
	return err
}

func (s *Storage) UserCookieGet(id int64) model.User {
	users := []model.User{}
	err := s.DB.Select(&users, "SELECT * FROM user WHERE id=?", id)

	if len(users) > 0 {
		return users[0]
	}
	if err != nil {
		fmt.Println(err)
	}
	return model.User{}
}

func (s *Storage) UserCookieDelete(user model.User) error {
	delete, err := s.DB.Prepare("DELETE FROM user WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	_, err = delete.Exec(user.ID)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
