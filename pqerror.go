package main

import (
	"fmt"

	"github.com/lib/pq"
)

func pqErrorDetector() {
	var err error
	if err, ok := err.(*pq.Error); ok {
		fmt.Println("Severity", err.Severity)
		fmt.Println("Code", err.Code)
		fmt.Println("Message", err.Message)
		fmt.Println("Detail", err.Detail)
		fmt.Println("Hint", err.Hint)
		fmt.Println("Position", err.Position)
		fmt.Println("InternalPosition", err.InternalPosition)
		fmt.Println("InternalQuery", err.InternalQuery)
		fmt.Println("Where", err.Where)
		fmt.Println("Schema", err.Schema)
		fmt.Println("Table", err.Table)
		fmt.Println("Column", err.Column)
		fmt.Println("DataTypeName", err.DataTypeName)
		fmt.Println("Constraint", err.Constraint)
		fmt.Println("File", err.File)
		fmt.Println("Line", err.Line)
		fmt.Println("Routine", err.Routine)
	}
	// Example
	// Severity ERROR
	// Code 23505
	// Message duplicate key value violates unique constraint "kuser_username_key"
	// Detail Key (username)=(rahulsoibam) already exists.
	// Hint
	// Position
	// InternalPosition
	// InternalQuery
	// Where
	// Schema public
	// Table kuser
	// Column
	// DataTypeName
	// Constraint kuser_username_key
	// File nbtinsert.c
	// Line 434
	// Routine _bt_check_unique
}
