//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/sqlite"
)

var Supertype = newSupertypeTable("", "Supertype", "")

type supertypeTable struct {
	sqlite.Table

	// Columns
	ID   sqlite.ColumnInteger
	Name sqlite.ColumnString

	AllColumns     sqlite.ColumnList
	MutableColumns sqlite.ColumnList
}

type SupertypeTable struct {
	supertypeTable

	EXCLUDED supertypeTable
}

// AS creates new SupertypeTable with assigned alias
func (a SupertypeTable) AS(alias string) *SupertypeTable {
	return newSupertypeTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new SupertypeTable with assigned schema name
func (a SupertypeTable) FromSchema(schemaName string) *SupertypeTable {
	return newSupertypeTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new SupertypeTable with assigned table prefix
func (a SupertypeTable) WithPrefix(prefix string) *SupertypeTable {
	return newSupertypeTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new SupertypeTable with assigned table suffix
func (a SupertypeTable) WithSuffix(suffix string) *SupertypeTable {
	return newSupertypeTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newSupertypeTable(schemaName, tableName, alias string) *SupertypeTable {
	return &SupertypeTable{
		supertypeTable: newSupertypeTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newSupertypeTableImpl("", "excluded", ""),
	}
}

func newSupertypeTableImpl(schemaName, tableName, alias string) supertypeTable {
	var (
		IDColumn       = sqlite.IntegerColumn("id")
		NameColumn     = sqlite.StringColumn("name")
		allColumns     = sqlite.ColumnList{IDColumn, NameColumn}
		mutableColumns = sqlite.ColumnList{NameColumn}
	)

	return supertypeTable{
		Table: sqlite.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:   IDColumn,
		Name: NameColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}