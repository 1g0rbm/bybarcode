package auth

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"-" db:"id"`
	Token        uuid.UUID `json:"token" db:"token"`
	RefreshToken uuid.UUID `json:"refresh_token" db:"refresh_token"`
	AccountID    int64     `json:"account_id" db:"account_id"`
	ExpireAt     time.Time `json:"expire_at" db:"expire_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (s Session) Encode() ([]byte, error) {
	return json.Marshal(s)
}
