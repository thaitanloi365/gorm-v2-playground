package main

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()

	var session = db.Session(&gorm.Session{PrepareStmt: true, DryRun: true})
	var statement = session.Statement

	statement.AddClause(clause.Where{
		Exprs: statement.BuildCondition("idemail = ?", "2"),
	})
	statement.AddClause(clause.Where{
		Exprs: statement.BuildCondition("id = ?", "2"),
	})

	statement.AddClause(clause.Limit{
		Limit:  10,
		Offset: 2,
	})
	statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprint("email desc"), Raw: true},
		}},
	})

	statement.Vars = append(statement.Vars, sql.Named("test", "1234"))

	statement.Build("SELECT", "FROM", "WHERE", "LIMIT", "ORDER BY")

	fmt.Println(statement.SQL.String())
	fmt.Println(statement.Vars...)

}
