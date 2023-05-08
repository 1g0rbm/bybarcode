package stat

import (
	"encoding/json"
	"io"
	"time"
)

type Statistic struct {
	ID                   int64     `json:"id"`
	ShoppingListId       int64     `json:"shopping_list_id"`
	ShoppingListName     string    `json:"shoppingListName"`
	CreatedAt            time.Time `json:"created_at"`
	ProductsCount        int       `json:"products_count"`
	CheckedProductsCount int       `json:"checked_products_count"`
}

func (s *Statistic) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Statistic) Decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return err
	}

	return nil
}
