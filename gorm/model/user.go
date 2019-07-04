package model

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// Create a GORM-backend model
type User struct {
	ID     uint `gorm:"primary_key"`
	Name   string
	SHA256 string
}

func (u *User) BeforeUpdate(scope *gorm.Scope) (err error) {
	json, err := json.Marshal(*u)
	sum := sha256.Sum256(json)
	sumB46 := base64.StdEncoding.EncodeToString(sum[:])
	scope.SetColumn("SHA256", sumB46)
	//fmt.Printf("Scope: %#v\n", *scope)
	return
}
