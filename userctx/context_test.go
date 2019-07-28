package userctx

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestUserCtx(t *testing.T) {
	ctx := context.Background()
	SayGreeting(ctx)
	ctx = WithUser(ctx, "cica")
	SayGreeting(ctx)
	value := ctx.Value("currentUser")
	fmt.Printf("Hello %v!\n", value)

}

func SayGreeting(ctx context.Context) {
	user := CurrentUser(ctx)
	fmt.Printf("Hello %s!\n", user)
}

type privateCtxKey string

const ctxKeyCurrentUser privateCtxKey = "currentUser"

func WithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, ctxKeyCurrentUser, user)
}

func CurrentUser(ctx context.Context) string {
	value := ctx.Value(ctxKeyCurrentUser)
	if value == nil {
		return ""
	}

	currentUser, isOk := value.(string)
	if isOk {
		return currentUser
	}

	panic(errors.Errorf("Isn't string: %v", value))
}
