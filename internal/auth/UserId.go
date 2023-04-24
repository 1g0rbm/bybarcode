package auth

import (
	"encoding/json"
	"io"
)

type UserId struct {
	Value int64 `json:"user_id"`
}

func (u *UserId) Decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&u); err != nil {
		return err
	}

	return nil
}
