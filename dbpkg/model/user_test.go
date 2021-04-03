package model

import (
	"context"
	"fmt"
	"github.com/bukodi/go-playgroud/dbpkg"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
	"time"
)

func TestUserInheritance(t *testing.T) {
	var u1 interface{} = &User{}
	{
		base, isOk := u1.(User)
		fmt.Printf("%v, %v\n", isOk, base)
	}
	{
		base, isOk := u1.(*User)
		fmt.Printf("%v, %v\n", isOk, base)
	}
	{
		base, isOk := u1.(dbpkg.ModelBase)
		fmt.Printf("%v, %v\n", isOk, base)
	}
	{
		base, isOk := u1.(*dbpkg.ModelBase)
		fmt.Printf("%v, %v\n", isOk, base)
	}
	{
		base, isOk := u1.(dbpkg.ModelIf)
		fmt.Printf("%v, %v\n", isOk, base)
	}
	{
		base, isOk := u1.(*dbpkg.ModelIf)
		fmt.Printf("%v, %v\n", isOk, base)
	}
}

func TestUserCRUD(t *testing.T) {

	t.Log("Starting test")

	os.Remove("test.db")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	db.LogMode(true)
	db.AutoMigrate(&User{})

	err = dbpkg.DoInTransaction(context.Background(), db, func(ctx context.Context) error {
		u1 := &User{
			Name: "Gipsz Jakab",
		}

		u1.ValidFrom = time.Now()

		if err := u1.Create(ctx); err != nil {
			t.Error(err)
		}
		t.Logf("User1 created: %#v", u1)

		tx := dbpkg.CurrentDB(ctx)
		var u2 User
		tx.First(&u2, "name = ?", "Gipsz Jakab")
		tx.Model(&u2).Update("Name", "Teszt Anna")

		users, err := AllUsers(ctx, 0, 0)
		if err != nil {
			t.Error(err)
		}
		t.Logf("List of users: : %#v", users)

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
