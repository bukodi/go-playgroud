package chirest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"testing"
)

type handlerWrapper struct {
	realHandler http.Handler
}

var _ http.Handler = &handlerWrapper{}

func (h handlerWrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.realHandler.ServeHTTP(writer, request)
}

func TestChiHello(t *testing.T) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	httpSrv, _ := StartServer(":3333", r)
	defer httpSrv.Shutdown(context.TODO())

	{
		resp, err := http.Get("http://localhost:3333/")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Response body: %s", string(body))
	}

}
