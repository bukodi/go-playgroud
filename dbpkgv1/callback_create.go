package dbpkgv1

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

// createCallback the callback used to insert data into database
func createCallback(scope *gorm.Scope, modelEntity ModelIf) {
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
