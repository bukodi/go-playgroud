package dbpkgv1

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type ModelBase struct {
	ID        uint      `sql:"index"`
	ValidFrom time.Time `sql:"index"`
	ValidTo   time.Time `sql:"index"`
	ModUserID uint
	PrevHash  string `sql:"index"`
	ThisHash  string `sql:"index,gorm:"primary_key"`
}

func (mb *ModelBase) AsModelBase() *ModelBase {
	return mb
}

func (mb *ModelBase) RecalcHash(now time.Time) (prevHash, newHash string) {
	prevHash = mb.PrevHash
	mb.PrevHash = mb.ThisHash
	mb.ThisHash = ""
	mb.ValidFrom = now
	mb.ValidTo = time.Time{}
	str := fmt.Sprintf("%+v", mb)
	sum := sha256.Sum256([]byte(str))
	newHash = base64.StdEncoding.EncodeToString(sum[:])
	mb.ThisHash = newHash
	return
}

func (mb *ModelBase) VerifyHash() error {
	savedThisHash := mb.ThisHash
	mb.ThisHash = ""
	savedValidTo := mb.ValidTo
	mb.ValidTo = time.Time{}

	defer func() {
		mb.ThisHash = savedThisHash
		mb.ValidTo = savedValidTo
	}()

	str := fmt.Sprintf("%+v", mb)
	sum := sha256.Sum256([]byte(str))
	calculatedHash := base64.StdEncoding.EncodeToString(sum[:])

	if savedThisHash != calculatedHash {
		return fmt.Errorf("Hash verification failed. Hash: %s, Record: %s", calculatedHash, str)
	}
	return nil
}

type ModelIf interface {
	AsModelBase() *ModelBase
	RecalcHash(time.Time) (string, string)
	VerifyHash() error
}
