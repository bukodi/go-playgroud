package dbpkgv1

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Create another GORM-backend model
type TestStructUser struct {
	ID     uint `gorm:"primary_key"`
	Name   string
	SHA256 string
}

type Address struct {
	ID     uint `gorm:"primary_key"`
	City   string
	Street string
}

func TestSingleEntity(t *testing.T) {
	os.Remove("test.db")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Address{})

	db.Callback().Update().Before("gorm:update").Register("logUpdate", logUpdate)
	db.Callback().Create().Before("gorm:create").Register("logUpdate", logUpdate)
	// Create
	db.Create(&Address{City: "Budapest"})

	// Read
	var addr Address
	db.First(&addr, 1) // find product with id 1
	fmt.Printf("Address: %#v\n", addr)
	db.First(&addr, "city = ?", "Budapest") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&addr).Update("City", "Kecskemét")

	// Delete - delete product
	//db.Delete(&user)
	//t.Error()
}

func TestTxUpdate(t *testing.T) {
	os.Remove("test.db")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(true)
	//db.SetLogger(gorm.Logger{revel.TRACE})
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	db.AutoMigrate(&TestStructUser{})
	db.Callback().Update().Before("gorm:update").Register("logUpdate", logUpdate)

	// Create
	db.Create(&TestStructUser{ID: 100, Name: "Teszt Anna"})
	db.Create(&TestStructUser{ID: 200, Name: "Teszt Bea"})

	tx := db.Begin()
	defer tx.Commit()
	// Migrate the schema

	// Read
	var user1 TestStructUser
	tx.First(&user1, "name = ?", "Teszt Anna")
	tx.Model(&user1).Update("Name", "Teszt András")

	var user2 TestStructUser
	tx.First(&user2, "name = ?", "Teszt Bea")
	tx.Model(&user2).Update("Name", "Teszt Béla")

	// Delete - delete product
	//db.Delete(&user)
	//t.Error()
}

func logUpdate(scope *gorm.Scope) {
	//scope.Commit()
	fmt.Printf("logUpdate: %+v\n", scope)
}
