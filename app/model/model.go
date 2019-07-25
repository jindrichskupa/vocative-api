package model

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	// import PostgreSQL dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// FirstName structure
type FirstName struct {
	Name       string `json:"name"`
	Vocative   string `json:"vocative"`
	NameSearch string `json:"name_search"`
	Count      uint   `json:"count"`
	Gender     string `json:"gender"`
}

// FirstNameJSON structure
type FirstNameJSON struct {
	Name     string `json:"name"`
	Vocative string `json:"vocative"`
	Count    uint   `json:"count"`
	Gender   string `json:"gender"`
}

// MarshalJSON converts FirstName to FirstNameJSON
func (n FirstName) MarshalJSON() ([]byte, error) {
	m := FirstNameJSON{
		Name:     n.Name,
		Vocative: n.Vocative,
		Count:    n.Count,
		Gender:   n.Gender,
	}
	return json.Marshal(m)
}

// SurName structure
type SurName struct {
	Name       string `json:"name"`
	Vocative   string `json:"vocative"`
	NameSearch string `json:"name_search"`
	Count      uint   `json:"count"`
	Gender     string `json:"gender"`
}

// SurNameJSON structure
type SurNameJSON struct {
	Name     string `json:"name"`
	Vocative string `json:"vocative"`
	Count    uint   `json:"count"`
	Gender   string `json:"gender"`
}

// MarshalJSON converts SurName to SurNameJSON
func (n SurName) MarshalJSON() ([]byte, error) {
	m := SurNameJSON{
		Name:     n.Name,
		Vocative: n.Vocative,
		Count:    n.Count,
		Gender:   n.Gender,
	}
	return json.Marshal(m)
}

// VocativeJSON structure
type VocativeJSON struct {
	Name     string `json:"name"`
	Vocative string `json:"vocative"`
	Count    uint   `json:"count"`
	Gender   string `json:"gender"`
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	return db
}
