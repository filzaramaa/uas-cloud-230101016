package models

import (
	"time"
)

type Pesan struct {
	Kode      string `gorm:"primaryKey;type:varchar(50)"`
	Balasan   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
