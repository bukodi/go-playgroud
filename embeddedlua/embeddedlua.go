package main

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

func main() {
	luaSrc := `
	output = string.reverse(input)
	`

	L := lua.NewState(lua.Options{
		CallStackSize: 10,
		SkipOpenLibs:  false,
	})
	L.SetGlobal("input", lua.LString("cica"))
	defer L.Close()
	if err := L.DoString(luaSrc); err != nil {
		panic(err)
	}
	lvOutput := L.GetGlobal("output")
	fmt.Println("Output: " + lvOutput.String())
}
