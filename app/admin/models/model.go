package models

import (
	"database/sql/driver"
	"encoding/json"
)

type BaseModel struct {
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt"`
}

type IntArray []int

func (j *IntArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok || len(bytes) == 0 {
		return nil
	}

	result := IntArray{}
	err := json.Unmarshal(bytes, &result)
	*j = IntArray(result)
	return err
}

func (j IntArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "", nil
	}

	bytes, err := json.Marshal(j)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}
