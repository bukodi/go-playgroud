package calldll

import (
	"fmt"
	"syscall"
	"testing"
	"unicode/utf16"
	"unsafe"
)

var (
	vscAgtDll = syscall.NewLazyDLL("c:\\tmp\\VSC\\NrgMrgVSCAgtUtils.dll")

	procGetProviderNames      = vscAgtDll.NewProc("GetProviderNames")
	procGetProvIndexedKeyProp = vscAgtDll.NewProc("GetProvIndexedKeyProp")
	//procGetProvNamedPubKey = vscAgtDll.NewProc("GetProvNamedPubKey")
)

func TestSyscall(t *testing.T) {
	name, _ := syscall.ComputerName()
	fmt.Println(name)
}
func TestCert(t *testing.T) {
	//TestSetProvider(t)
}

const MAX_LENGTH = 4095

func TestListProviders(t *testing.T) {

	//TestSetProvider(t)

	var n uint32 = MAX_LENGTH + 1
	b := make([]uint16, n)

	r1, _, e1 := syscall.Syscall(procGetProviderNames.Addr(), 2, uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&n)), 0)
	if r1 != 0 {
		fmt.Printf("r1=%v, e1=%v\n", r1, e1)
	}

	var providerNames = string(utf16.Decode(b[0 : n-1]))
	fmt.Printf("providerNames=%v\n", providerNames)

}

func TestListKeys(t *testing.T) {

	var pszParamList = "ProvName=Microsoft Software Key Storage Provider"
	pszParamList = pszParamList + ""

	var n uint32 = MAX_LENGTH + 1
	b := make([]uint16, n)

	var idx uint16 = 10
	var nkpType uint16 = 1

	r1, r2, e1 := syscall.Syscall6(procGetProvIndexedKeyProp.Addr(), 5,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pszParamList))),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(&n)),
		uintptr(idx),
		uintptr(nkpType),
		0)
	if r1 != 0 {
		fmt.Printf("r1=%x, r2=%v, e1=%v\n", r1, r2, e1)
	} else {
		var keyproperty = string(utf16.Decode(b[0 : n-1]))
		fmt.Printf("%d key name=%v\n", idx, keyproperty)
	}

}
