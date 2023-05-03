package listener

import (
	"bybarcode/internal/db"
	"context"
)

type EventListener struct {
	db       *db.Connect
	slIdChan chan int64
}

func NewEventListener(db *db.Connect) *EventListener {
	return &EventListener{
		db:       db,
		slIdChan: make(chan int64),
	}
}

func (e *EventListener) Listen(ctx context.Context) error {
	for {
		select {
		case slId := <-e.slIdChan:
			if err := e.db.AddedUpdStatisticByShoppingList(ctx, slId); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (e *EventListener) Notify(slId int64) {
	go func() {
		e.slIdChan <- slId
	}()
}
