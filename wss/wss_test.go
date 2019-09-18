package wss

import (
	"github.com/sacOO7/gowebsocket"
	"log"
	"os"
	"os/signal"
	"testing"
)

func TestWSSCall(t *testing.T) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("ws://172.16.149.1:36080/mdaserver/mda");

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server");
	};

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	};

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	};

	socket.OnBinaryMessage = func(data [] byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	};

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	};

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	};

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	};

	socket.Connect()

	socket.SendText(`{"method":"mdaConfig","requestId":1,"headers":{"X-GSSToken":"NagyonTesztToken","X-DeviceId":"enyime"},"data":[]}`)

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}
