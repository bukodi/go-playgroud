package mybolt

import (
	"log"

	bolt "github.com/coreos/bbolt"
	"testing"
	"fmt"
)

func TestDb(t *testing.T) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("/tmp/mybolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println(db.Path())

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("Cica"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		b.Put([]byte("key01"), []byte("value01"))
		return nil
	})
}