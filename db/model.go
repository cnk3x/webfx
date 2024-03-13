package db

import (
	"time"

	"gorm.io/gorm"
)

type Tabled interface {
	GetPK() int
	GetCreatedAt() int
}

type Table struct {
	ID        int `json:"id"`
	CreatedAt int `json:"created_at"`
}

func (m Table) GetPK() int {
	return m.ID
}

func (m Table) GetCreatedAt() int {
	return m.CreatedAt
}

func (m *Table) BeforeCreate(tx *gorm.DB) (err error) {
	if m.CreatedAt == 0 {
		m.CreatedAt = HumanDate(time.Now())
	}
	return
}

func ID(id int) Table {
	return Table{ID: id}
}
