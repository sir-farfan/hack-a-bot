package sql_event

import (
	"fmt"
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
		fmt.Printf("deleting: %v\n", evt)
		err = db.DeleteEvent(evt.ID)
		assert.Nil(t, err)
	}

	db.CreateEvent(model.Event{Name: "silks", Description: "aerial silks dance"})
	db.CreateEvent(model.Event{Name: "straps", Description: "aerial straps skills"})
	db.CreateEvent(model.Event{Name: "pole", Description: "pole fitness"})
}
