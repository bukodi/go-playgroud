package main

import(
	"github.com/ethereum/go-ethereum/ethdb"
	"fmt"
)

func main() {
	db, err := ethdb.NewLDBDatabase("/home/bukodi/.ethereum/testnet/geth/chaindata", 100, 100)
	if err != nil {
		panic("failed to create test database: " + err.Error())
	}
	it := db.NewIterator()
	fmt.Println()

}
