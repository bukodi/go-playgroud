package dbpkg

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type privateCtxKey string

const ctxKeyCurrentDb privateCtxKey = "currentDb"

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyCurrentDb, db)
}

func CurrentDB(ctx context.Context) *gorm.DB {
	value := ctx.Value(ctxKeyCurrentDb)
	if value == nil {
		return nil
	}

	currentDB, isOk := value.(*gorm.DB)
	if isOk {
		return currentDB
	}

	panic(errors.Errorf("Isn't a grom.DB: %v", value))
}

type InTransaction func(ctx context.Context) error

func DoInTransaction(ctx context.Context, db *gorm.DB, fn InTransaction) error {
	tx := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
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
		modelIf, isOk := scope.Value.(ModelIf)
		if !isOk {
			originalCreateProcessor(scope)
		} else {
			fmt.Printf("ModelBase: %v\n", modelIf.AsModelBase())
			createCallback(scope)
		}
	})

	tx.Callback().Update().Remove("gorm:begin_transaction")
	tx.Callback().Update().Remove("gorm:commit_or_rollback_transaction")
	tx.Callback().Update().Before("gorm:update").Register("txLogUpdate", func(scope *gorm.Scope) {
		logTxOperation(scope, "UPDATE")
	})
	originalUpdateProcessor := tx.Callback().Update().Get("gorm:update")
	tx.Callback().Update().Replace("gorm:update", func(scope *gorm.Scope) {
		modelIf, isOk := scope.Value.(ModelIf)
		if !isOk {
			originalUpdateProcessor(scope)
		} else {
			fmt.Printf("ModelBase: %v\n", modelIf.AsModelBase())
			updateCallback(scope)
		}
	})

	tx.Callback().Delete().Remove("gorm:begin_transaction")
	tx.Callback().Delete().Remove("gorm:commit_or_rollback_transaction")
	tx.Callback().Delete().Before("gorm:delete").Register("txLogDelete", func(scope *gorm.Scope) {
		logTxOperation(scope, "DELETE")
	})

	// Store tx in the context
	dbCtx := context.WithValue(ctx, ctxKeyCurrentDb, db)

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
	fmt.Printf("Tx %s: %#v\n", operation, scope.Value)
}

type TxLogEntry struct {
	Operation string
	Table     string
	Hash      string
	PrevHash  string
}

// updateCallback the callback used to update data to database
func updateCallback(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}
	var sqls []string

	if updateAttrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
		// Sort the column names so that the generated SQL is the same every time.
		updateMap := updateAttrs.(map[string]interface{})
		var columns []string
		for c := range updateMap {
			columns = append(columns, c)
		}
		sort.Strings(columns)

		for _, column := range columns {
			value := updateMap[column]
			sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(column), scope.AddToVars(value)))
		}
	} else {
		for _, field := range scope.Fields() {
			if changeableFieldOfScope(scope, field) {
				if !field.IsPrimaryKey && field.IsNormal && (field.Name != "CreatedAt" || !field.IsBlank) {
					if !field.IsForeignKey || !field.IsBlank || !field.HasDefaultValue {
						sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(field.DBName), scope.AddToVars(field.Field.Interface())))
					}
				} else if relationship := field.Relationship; relationship != nil && relationship.Kind == "belongs_to" {
					for _, foreignKey := range relationship.ForeignDBNames {
						if foreignField, ok := scope.FieldByName(foreignKey); ok && !changeableFieldOfScope(scope, foreignField) {
							sqls = append(sqls,
								fmt.Sprintf("%v = %v", scope.Quote(foreignField.DBName), scope.AddToVars(foreignField.Field.Interface())))
						}
					}
				}
			}
		}
	}

	var extraOption string
	if str, ok := scope.Get("gorm:update_option"); ok {
		extraOption = fmt.Sprint(str)
	}

	if len(sqls) > 0 {
		scope.Raw(fmt.Sprintf(
			"UPDATE %v SET %v%v%v",
			scope.QuotedTableName(),
			strings.Join(sqls, ", "),
			addExtraSpaceIfExist(scope.CombinedConditionSql()),
			addExtraSpaceIfExist(extraOption),
		)).Exec()
	}
}

// createCallback the callback used to insert data into database
func createCallback(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}
	//defer scope.trace(scope.db.nowFunc())

	var (
		columns, placeholders        []string
		blankColumnsWithDefaultValue []string
	)

	for _, field := range scope.Fields() {
		if changeableFieldOfScope(scope, field) {
			if field.IsNormal && !field.IsIgnored {
				if field.IsBlank && field.HasDefaultValue {
					blankColumnsWithDefaultValue = append(blankColumnsWithDefaultValue, scope.Quote(field.DBName))
					scope.InstanceSet("gorm:blank_columns_with_default_value", blankColumnsWithDefaultValue)
				} else if !field.IsPrimaryKey || !field.IsBlank {
					columns = append(columns, scope.Quote(field.DBName))
					placeholders = append(placeholders, scope.AddToVars(field.Field.Interface()))
				}
			} else if field.Relationship != nil && field.Relationship.Kind == "belongs_to" {
				for _, foreignKey := range field.Relationship.ForeignDBNames {
					if foreignField, ok := scope.FieldByName(foreignKey); ok && !changeableFieldOfScope(scope, foreignField) {
						columns = append(columns, scope.Quote(foreignField.DBName))
						placeholders = append(placeholders, scope.AddToVars(foreignField.Field.Interface()))
					}
				}
			}
		}
	}

	var (
		returningColumn = "*"
		quotedTableName = scope.QuotedTableName()
		primaryField    = scope.PrimaryField()
		extraOption     string
		insertModifier  string
	)

	if str, ok := scope.Get("gorm:insert_option"); ok {
		extraOption = fmt.Sprint(str)
	}
	if str, ok := scope.Get("gorm:insert_modifier"); ok {
		insertModifier = strings.ToUpper(fmt.Sprint(str))
		if insertModifier == "INTO" {
			insertModifier = ""
		}
	}

	if primaryField != nil {
		returningColumn = scope.Quote(primaryField.DBName)
	}

	lastInsertIDReturningSuffix := scope.Dialect().LastInsertIDReturningSuffix(quotedTableName, returningColumn)

	if len(columns) == 0 {
		scope.Raw(fmt.Sprintf(
			"INSERT %v INTO %v %v%v%v",
			addExtraSpaceIfExist(insertModifier),
			quotedTableName,
			scope.Dialect().DefaultValueStr(),
			addExtraSpaceIfExist(extraOption),
			addExtraSpaceIfExist(lastInsertIDReturningSuffix),
		))
	} else {
		scope.Raw(fmt.Sprintf(
			"INSERT %v INTO %v (%v) VALUES (%v)%v%v",
			addExtraSpaceIfExist(insertModifier),
			scope.QuotedTableName(),
			strings.Join(columns, ","),
			strings.Join(placeholders, ","),
			addExtraSpaceIfExist(extraOption),
			addExtraSpaceIfExist(lastInsertIDReturningSuffix),
		))
	}

	// execute create sql
	if lastInsertIDReturningSuffix == "" || primaryField == nil {
		if result, err := scope.SQLDB().Exec(scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
			// set rows affected count
			scope.DB().RowsAffected, _ = result.RowsAffected()

			// set primary value to primary field
			if primaryField != nil && primaryField.IsBlank {
				if primaryValue, err := result.LastInsertId(); scope.Err(err) == nil {
					scope.Err(primaryField.Set(primaryValue))
				}
			}
		}
	} else {
		if primaryField.Field.CanAddr() {
			if err := scope.SQLDB().QueryRow(scope.SQL, scope.SQLVars...).Scan(primaryField.Field.Addr().Interface()); scope.Err(err) == nil {
				primaryField.IsBlank = false
				scope.DB().RowsAffected = 1
			}
		} else {
			scope.Err(gorm.ErrUnaddressable)
		}
	}
}

func changeableFieldOfScope(scope *gorm.Scope, field *gorm.Field) bool {
	if selectAttrs := scope.SelectAttrs(); len(selectAttrs) > 0 {
		for _, attr := range selectAttrs {
			if field.Name == attr || field.DBName == attr {
				return true
			}
		}
		return false
	}

	for _, attr := range scope.OmitAttrs() {
		if field.Name == attr || field.DBName == attr {
			return false
		}
	}

	return true
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
