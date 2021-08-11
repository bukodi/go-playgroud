package dbpkgv1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type privateCtxKey string

const ctxKey privateCtxKey = "txInfo"

type txInfo struct {
	currentDB *gorm.DB
	txTime    time.Time
	log       []TxLogEntry
}

type TxLogEntry struct {
	Operation string
	Table     string
	Hash      string
	UpdateSQL string
	CreateSQL string
}

func CurrentDB(ctx context.Context) *gorm.DB {
	value := ctx.Value(ctxKey)
	if value == nil {
		return nil
	}

	ti, isOk := value.(*txInfo)
	if isOk {
		return ti.currentDB
	}

	panic(errors.Errorf("Isn't a grom.DB: %v", value))
}

type InTransaction func(ctx context.Context) error

func DoInTransaction(ctx context.Context, db *gorm.DB, fn InTransaction) error {
	tx := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if tx.Error != nil {
		return tx.Error
	}

	// Use same timestamp in the transaction
	txTimestamp := time.Now().UTC()
	tx.SetNowFuncOverride(func() time.Time {
		return txTimestamp
	})

	// Add callbacks
	tx.Callback().Create().Remove("gorm:begin_transaction")
	tx.Callback().Create().Remove("gorm:commit_or_rollback_transaction")
	tx.Callback().Create().Before("gorm:create").Register("txLogCreate", func(scope *gorm.Scope) {
		logTxOperation(scope, "CREATE")
	})
	originalCreateProcessor := tx.Callback().Create().Get("gorm:create")
	tx.Callback().Create().Replace("gorm:create", func(scope *gorm.Scope) {
		modelEntity, isOk := scope.Value.(ModelIf)
		if !isOk {
			originalCreateProcessor(scope)
		} else {
			createCallback(scope, modelEntity)
		}
	})

	tx.Callback().Update().Remove("gorm:begin_transaction")
	tx.Callback().Update().Remove("gorm:commit_or_rollback_transaction")
	tx.Callback().Update().Before("gorm:update").Register("txLogUpdate", func(scope *gorm.Scope) {
		logTxOperation(scope, "UPDATE")
	})
	originalUpdateProcessor := tx.Callback().Update().Get("gorm:update")
	tx.Callback().Update().Replace("gorm:update", func(scope *gorm.Scope) {
		modelEntity, isOk := scope.Value.(ModelIf)
		if !isOk {
			originalUpdateProcessor(scope)
		} else {
			updateCallback(scope, modelEntity)
		}
	})

	tx.Callback().Delete().Remove("gorm:begin_transaction")
	tx.Callback().Delete().Remove("gorm:commit_or_rollback_transaction")
	tx.Callback().Delete().Before("gorm:delete").Register("txLogDelete", func(scope *gorm.Scope) {
		logTxOperation(scope, "DELETE")
	})

	// Store tx in the context
	ti := txInfo{
		currentDB: tx,
		txTime:    txTimestamp,
		log:       make([]TxLogEntry, 0),
	}
	tx.InstantSet("TxInfo", &ti)
	dbCtx := context.WithValue(ctx, ctxKey, &ti)

	err := fn(dbCtx)
	if err != nil {
		xerr := tx.Rollback().Error
		if xerr != nil {
			return xerr
		}
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func logTxOperation(scope *gorm.Scope, operation string) {
	//scope.Commit()
	//fmt.Printf("Tx %s: %#v\n", operation, scope.Value)
	scope.Log(fmt.Sprintf("Tx %s: %#v\n", operation, scope.Value))

}
