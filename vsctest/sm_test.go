package vsctest

import (
	"fmt"
	"github.com/sf1/go-card/smartcard"
	"testing"
)

func TestVSC(t *testing.T) {
	ctx, err := smartcard.EstablishContext()
	if err != nil {
		t.Error(err)
	}
	// handle error, if any
	defer ctx.Release()

	readers, err := ctx.ListReadersWithCard()
	if err != nil { t.Error(err); return }
	for _, reader := range readers {
		fmt.Println(reader.Name())
		fmt.Printf("- Card present: %t\n\n", reader.IsCardPresent())
	}

	reader := readers[0]
	fmt.Println("Connect to card")
	fmt.Printf("---------------\n\n")
	card, err := reader.Connect()
	if err != nil { t.Error(err); return }
	fmt.Printf("ATR: %s\n\n", card.ATR())

	fmt.Println("Select applet")
	fmt.Printf("-------------\n\n")
	cmd := smartcard.SelectCommand(0x90, 0x72, 0x5A, 0x9E, 0x3B, 0x10, 0x70, 0xAA)
	fmt.Printf(">> %s\n", cmd)
	response, err := card.TransmitAPDU(cmd)
	if err != nil { t.Error(err); return }
	fmt.Printf("<< %s\n", response.String())

	fmt.Println("\nSend CMD 10")
	fmt.Printf("-----------\n\n")
	cmd = smartcard.Command2(0x00, 0x10, 0x00, 0x00, 0x0b)
	fmt.Printf(">> %s\n", cmd)
	response, err = card.TransmitAPDU(cmd)
	if err != nil { t.Error(err); return }
	fmt.Printf("<< %s\n", response)
	fmt.Printf("\nQuoth the Applet, \"%s\"\n\n", string(response.Data()))

	fmt.Println("\nSend CHANGE PIN CMD")
	fmt.Printf("-----------\n\n")
	cmd = smartcard.Command2(0x00, 0x10, 0x00, 0x00, 0x0b)
	fmt.Printf(">> %s\n", cmd)
	response, err = card.TransmitAPDU(cmd)
	if err != nil { t.Error(err); return }
	fmt.Printf("<< %s\n", response)
	fmt.Printf("\nQuoth the Applet, \"%s\"\n\n", string(response.Data()))

	fmt.Println("Disconnect from card")
	fmt.Printf("--------------------\n\n")
	err = card.Disconnect()
	if err != nil { t.Error(err); return }
	fmt.Printf("OK\n\n")
}

