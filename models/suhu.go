package models

import (
	"time"
)

type Suhu struct {
	Id        uint      `gorm:"primary_key"`
	Lokasi    string    `gorm:"type:varchar(225)"`
	Suhu      float32   `gorm:"type:decimal(10,2)"`
	CreatedAt time.Time `gorm:"type:datetime"`
}
