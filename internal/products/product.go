package products

import (
	"encoding/json"
	"fmt"
	"io"
)

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Brand struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Upcean     int64     `json:"upcean"`
	CategoryId int64     `json:"category_id"`
	BrandId    int64     `json:"brand_id"`
	Category   *Category `json:"category,omitempty"`
	Brand      *Brand    `json:"brand,omitempty"`
}

type ProductInList struct {
	Product
	Checked bool `json:"checked,omitempty"`
}

func (p *ProductInList) MarshalJSON() ([]byte, error) {
	type Alias ProductInList
	fmt.Println(p.Checked)
	if p.Checked {
		return json.Marshal(&struct {
			Alias
			Checked bool `json:"checked"`
		}{
			Alias:   (Alias)(*p),
			Checked: p.Checked,
		})
	} else {
		return json.Marshal(&struct {
			Alias
			Checked bool `json:"checked"`
		}{
			Alias:   (Alias)(*p),
			Checked: p.Checked,
		})
	}
}

func (p *Product) Encode() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Product) Decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return err
	}

	return nil
}
