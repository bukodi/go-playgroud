package model

import (
	"context"

	"github.com/bukodi/go-playgroud/dbpkg"
)

// Create a GORM-backend model
type User struct {
	dbpkg.ModelBase
	Name   string
	SHA256 string
}

/*
func (u *User) BeforeUpdate(scope *gorm.Scope) (err error) {
	json, err := json.Marshal(*u)
	sum := sha256.Sum256(json)
	sumB46 := base64.StdEncoding.EncodeToString(sum[:])
	scope.SetColumn("SHA256", sumB46)
	//fmt.Printf("Scope: %#v\n", *scope)
	return
}*/

func (u *User) Create(ctx context.Context) error {
	tx := dbpkg.CurrentDB(ctx)
	tx.Create(u)
	return tx.Error
}

func (u *User) Update(ctx context.Context) error {
	tx := dbpkg.CurrentDB(ctx)
	tx.Save(u)
	return tx.Error
}

func AllUsers(ctx context.Context, firstId uint, limit int) ([]User, error) {
	var users []User
	tx := dbpkg.CurrentDB(ctx)
	tx.Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}
