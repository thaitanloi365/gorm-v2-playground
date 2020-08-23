package query

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func addClause(statement *gorm.Statement, clauses map[string]clause.Clause) (addedClauses []string) {
	addedClauses = []string{}

	for key, value := range clauses {
		switch key {
		case "SELECT":
			statement.AddClause(value.Expression.(clause.Select))
			addedClauses = append(addedClauses, key)
		case "FROM":
			statement.AddClause(value.Expression.(clause.From))
			addedClauses = append(addedClauses, key)
		case "WHERE":
			statement.AddClause(value.Expression.(clause.Where))
			fmt.Println("Exprs", value.Expression.(clause.Where).Exprs)
			addedClauses = append(addedClauses, key)
		case "GROUP BY":
			statement.AddClause(value.Expression.(clause.GroupBy))
			addedClauses = append(addedClauses, key)
		case "ORDER BY":
			statement.AddClause(value.Expression.(clause.OrderBy))
			addedClauses = append(addedClauses, key)
		case "LIMIT":
			statement.AddClause(value.Expression.(clause.Limit))
			addedClauses = append(addedClauses, key)
			fmt.Println("Exprs", value.Expression.(clause.Limit))

		}

	}

	return
}
