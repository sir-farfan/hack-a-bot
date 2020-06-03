package sql_event

import (
	"log"
	"testing"

	"github.com/sir-farfan/hack-a-bot/model"
	"github.com/stretchr/testify/assert"
)

var testEvent = model.Event{Name: "test event CRUD", Description: "making sure sql works"}

func TestEventCreate(t *testing.T) {
	db := New()
	defer db.DB.Close()

	err := db.CreateEvent(testEvent)
	assert.Nil(t, err)

	events := db.GetEvent("")
	log.Println(events)
	assert.NotEmpty(t, events)

	for _, evt := range events {
		err = db.DeleteEvent(evt.ID)
		assert.Nil(t, err)
	}
}
