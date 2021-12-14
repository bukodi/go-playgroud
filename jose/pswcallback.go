package jose

import (
	"context"
)

type PswCallback func(hint string) []byte

type CacheablePswCallback func(hint string) (password []byte, pswCtx context.Context)

func NewPswCallback(password string) PswCallback {
	return func(hint string) []byte {
		return []byte(password)
	}
}
