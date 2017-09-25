package main

import (
	//	"github.com/syndtr/goleveldb/leveldb/iterator"
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
)

func main() {
	db, err := ethdb.NewLDBDatabase("/home/bukodi/.ethereum/testnet/geth/chaindata", 100, 100)
	if err != nil {
		panic("failed to create test database: " + err.Error())
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()
	genesis := core.Genesis{Config: params.AllProtocolChanges, Alloc: alloc}
	blockchain = core.NewBlockChain(db,)

	it := db.NewIterator()
	defer func() {
		if it != nil {
			it.Release()
		}
	}()

	const MAX int = 100

	fmt.Println("--- BEGIN ---")
	for cnt := 0; it.Next() && cnt < MAX; cnt++ {
		key := it.Key()
		fmt.Println(cnt, key)
	}
	fmt.Println("--- END ---")

}
