package models

import (
	"encoding/json"
)

const DefaultStatus = "TO DO"

// ToDo defines to-do item structure.
type ToDo struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

// FromJson creates ToDo object from JSON byte array.
func FromJson(data []byte) (ToDo, error) {
	var item ToDo
	err := json.Unmarshal(data, &item)
	if err != nil {
		return item, err
	}
	return item, nil
}
