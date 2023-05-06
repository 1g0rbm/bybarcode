package listener

import (
	"bybarcode/internal/db"
	"context"
	"time"
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

func (e *EventListener) Notify(ctx context.Context, slId int64) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	select {
	case e.slIdChan <- slId:
	case <-ctx.Done():
		return
	}
}
