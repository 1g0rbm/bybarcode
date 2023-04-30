package products

import (
	"bybarcode/internal/auth"
	"encoding/json"
	"io"
)

type ShoppingList struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	AccountId int64         `json:"account_id"`
	Account   *auth.Account `json:"account,omitempty"`
	Products  []Product     `json:"products,omitempty"`
}

func (sl *ShoppingList) Encode() ([]byte, error) {
	return json.Marshal(sl)
}

func (sl *ShoppingList) Decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&sl); err != nil {
		return err
	}

	return nil
}
