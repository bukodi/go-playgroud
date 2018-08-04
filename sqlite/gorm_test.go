package main

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Table1 struct {
	gorm.Model
	Name string `gorm:"size:255"` // Default size for string is 255, reset it with this tag
}

func TestGORMSqlite(t *testing.T) {
	fmt.Println("Start")
	db, _ := gorm.Open("sqlite3", "/tmp/gorm2.db")
	//db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	x := db.AutoMigrate(&Table1{})

	//row1 := Table1{Name: "Kica"}
	//db.Create(&row1)
	var rows []Table1
	db.Where("name LIKE ?", "%mi%").Find(&rows)
	for i, r := range rows {
		fmt.Printf("%d: %q\n", i, r.Name)
	}
	fmt.Println(x)
}
