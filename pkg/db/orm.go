package db

import (
	"chocolate/internal/firstAdmin"
	"chocolate/internal/passwords"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"fmt"
	"slices"
	"strconv"
)

var DB *gorm.DB

// OpenDB инициализирует базу данных
func OpenDB(dbName string) {
	DB, _ = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	DB.AutoMigrate(&Category{}, &Admin{}, &Item{})
}
func FillDefaultData(defaultCategories []string) {
	flag := true

	existedCategories := FindAllCategories()
	for i, defcategory := range defaultCategories {
		isFound := slices.Contains(existedCategories, defcategory)
		if !isFound {
			CreateCategory(defaultCategories[i])
			flag = false
		}
	}

	existedAdmins := FindAllAdmins()
	isFound := slices.Contains(existedAdmins, firstAdmin.Name)
	if !isFound {
		CreateAdmin(firstAdmin.Name, firstAdmin.Password, false)
		flag = false
	}

	if flag {
		fmt.Println("All default admins and categories are existed")
	} else {
		fmt.Println("All missing categories and admins added")
	}
}
func CreateCategory(name string) {
	tx := DB.Begin()

	tx.Create(&Category{Name: name})

	tx.Commit()
}

// если hash == true, то пароль надо захешировать
func CreateAdmin(name string, password string, hash_pass bool) {
	if hash_pass {
		password = passwords.Hash(password)
	}
	tx := DB.Begin()

	tx.Create(&Admin{Name: name, Password: password})

	tx.Commit()
}
func CreateItem(name, longdesc, shortdesc string, price uint, category Category) {
	tx := DB.Begin()

	tx.Create(&Item{Name: name, Price: price, LongDesc: longdesc, ShortDesc: shortdesc, Category: category})

	tx.Commit()
}
func FindCategoryItems(category string) []Item {
	tx := DB.Begin()

	var items []Item
	tx.Where("category_name = ?", category).Find(&items)
	tx.Commit()
	return items
}
func AddItems(categories []string) {
	for i, category := range categories {
		strI := strconv.Itoa(i)
		CreateItem("item"+strI, "some long description"+strI, "short desc"+strI, uint(i), Category{Name: category})
	}
	fmt.Println("items from categories had created")
}
func FindAdminsPasswords(name string) []string {
	var admins []Admin

	tx := DB.Begin()

	tx.Where("name = ?", name).Find(&admins)

	tx.Commit()

	var result []string
	for _, admin := range admins {
		result = append(result, admin.Password)
	}
	return result
}
func FindAllCategories() []string {
	tx := DB.Begin()

	var categories []Category
	var result []string
	tx.Find(&categories)

	tx.Commit()
	for _, category := range categories {
		result = append(result, category.Name)
	}
	return result
}
func FindItems(name string) []Item {
	tx := DB.Begin()

	var items []Item
	tx.Where("UPPER(name) LIKE UPPER(? || '%')", name).Find(&items)

	tx.Commit()
	return items
}
func DeleteItem(name, category string) {
	tx := DB.Begin()

	tx.Where("UPPER(name) = UPPER(?) AND category_name = ?", name, category).Delete(&Item{})

	tx.Commit()
}
func DeleteCategory(name string) {
	tx := DB.Begin()

	tx.Where("UPPER(name) = UPPER(?)", name).Delete(&Category{})

	tx.Commit()
}
func FindAllAdmins() []string {
	tx := DB.Begin()

	var admins []Admin
	var result []string
	tx.Find(&admins)

	tx.Commit()
	for _, admin := range admins {
		result = append(result, admin.Name)
	}
	return result
}
