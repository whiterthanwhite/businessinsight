package category

import (
	"encoding/json"
	"log"

	"github.com/whiterthanwhite/businessinsight/internal/entities/operation_type"
)

type Category struct {
	Id          int                          `json:"id"`
	Type        operation_type.OperationType `json:"type"`
	Name        string                       `json:"name"`
	Description string                       `json:"description,omitempty"`
}

func (c *Category) UnmarshalJSON(body []byte) error {
	type temp struct {
		Id          int    `json:"id"`
		Type        string `json:"type"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}

	var t temp
	err := json.Unmarshal(body, &t)
	if err != nil {
		log.Println(string(body))
		return err
	}

	c.Id = t.Id
	c.Type = operation_type.OperationType(t.Type)
	c.Name = t.Name
	c.Description = t.Description

	return nil
}

func ParseJSON(categoriesJSON []byte) ([]Category, error) {
	var categories []Category
	err := json.Unmarshal(categoriesJSON, &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
