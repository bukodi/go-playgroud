package rpc

import (
	"net/rpc"
	"testing"
)


func TestRPC(t *testing.T) {
	listener :=

	srv := rpc.NewServer()
	srv.Register( impl )
	srv.Accept()

}
