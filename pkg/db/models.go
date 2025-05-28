package db

import (
	_ "gorm.io/gorm"
)

type Category struct {
	Name string
}
type Item struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Price     uint
	LongDesc  string
	ShortDesc string
	Category  Category `gorm:"embedded;embeddedPrefix:category_"`
}
type Admin struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Password string
}
