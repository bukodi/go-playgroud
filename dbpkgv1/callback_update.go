package dbpkgv1

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"sort"
	"strings"
)

// updateCallback the callback used to update data to database
func updateCallback(scope *gorm.Scope, modelEntity ModelIf) {
	if scope.HasError() {
		return
	}

	value, ok := scope.Get("TxInfo")
	if !ok {
		scope.Err(fmt.Errorf("Can't accuire TxInfo"))
		return
	}

	txi, ok2 := value.(*txInfo)
	if !ok2 {
		scope.Err(fmt.Errorf("Can't accuire TxInfo"))
		return
	}
	fmt.Printf("txInfo: : %+v\n ", txi)

	prevHash, newHash := modelEntity.RecalcHash(txi.txTime)
	fmt.Printf("Update %s -> %s\n", prevHash, newHash)

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
		sql := fmt.Sprintf(
			"UPDATE %v SET %v%v%v",
			scope.QuotedTableName(),
			strings.Join(sqls, ", "),
			addExtraSpaceIfExist(scope.CombinedConditionSql()),
			addExtraSpaceIfExist(extraOption),
		)
		fmt.Printf("Execute update: %s\n%+v", sql, scope.SQLVars)
		scope.Raw(sql).Exec()

	}
}
