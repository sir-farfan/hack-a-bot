package sql_event

import (
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

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS event (id INTEGER PRIMARY KEY, name TEXT, description TEXT)")
	_, err = statement.Exec()
	if err != nil {
		log.Printf("error creating database: %v\n", err)
	}

	return &Storage{DB: db}
}

func (s *Storage) GetEvent(name string) []model.Event {
	events := []model.Event{}
	err := s.DB.Select(&events, "SELECT id, name, description FROM event")
	if err != nil {
		log.Println(err)
	}

	return events
}

func (s *Storage) CreateEvent(event model.Event) error {
	insert, err := s.DB.Prepare("INSERT INTO event (name, description) VALUES (?, ?)")
	if err != nil {
		log.Printf("Error during insert event: %v\n", err)
		return err
	}
	_, err = insert.Exec(event.Name, event.Description)

	return err
}

func (s *Storage) DeleteEvent(id int) error {
	delete, err := s.DB.Prepare("DELETE FROM event WHERE id=?")
	if err != nil {
		return err
	}
	_, err = delete.Exec(id)
	return err
}
