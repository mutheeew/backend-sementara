package models

import "time"

type Category struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Name      string    `json:"name" form:"name" gorm:"type: varchar(255)"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name" form:"name" gorm:"type: varchar(255)"`
}

func (CategoryResponse) TableName() string {
	return "categories"
}
