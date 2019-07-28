package dbpkg

import "time"

type ModelBase struct {
	ID        uint      `gorm:"primary_key"`
	ValidFrom time.Time `sql:"index"`
	ValidTo   time.Time `sql:"index"`
	ModUserID uint
	PrevHash  string `sql:"index"`
	ThisHash  string `sql:"index"`
}

func (mb *ModelBase) AsModelBase() *ModelBase {
	return mb
}

type ModelIf interface {
	AsModelBase() *ModelBase
}
