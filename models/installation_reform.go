package models

// generated with gopkg.in/reform.v1

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type installationTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *installationTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("installations").
func (v *installationTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *installationTableType) Columns() []string {
	return []string{"id", "username", "source", "installation_id", "created_at", "updated_at"}
}

// NewStruct makes a new struct for that view or table.
func (v *installationTableType) NewStruct() reform.Struct {
	return new(Installation)
}

// NewRecord makes a new record for that table.
func (v *installationTableType) NewRecord() reform.Record {
	return new(Installation)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *installationTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// InstallationTable represents installations view or table in SQL database.
var InstallationTable = &installationTableType{
	s: parse.StructInfo{Type: "Installation", SQLSchema: "", SQLName: "installations", Fields: []parse.FieldInfo{{Name: "ID", PKType: "int64", Column: "id"}, {Name: "Username", PKType: "", Column: "username"}, {Name: "Source", PKType: "", Column: "source"}, {Name: "InstallationID", PKType: "", Column: "installation_id"}, {Name: "CreatedAt", PKType: "", Column: "created_at"}, {Name: "UpdatedAt", PKType: "", Column: "updated_at"}}, PKFieldIndex: 0},
	z: new(Installation).Values(),
}

// String returns a string representation of this struct or record.
func (s Installation) String() string {
	res := make([]string, 6)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Username: " + reform.Inspect(s.Username, true)
	res[2] = "Source: " + reform.Inspect(s.Source, true)
	res[3] = "InstallationID: " + reform.Inspect(s.InstallationID, true)
	res[4] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[5] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Installation) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Username,
		s.Source,
		s.InstallationID,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Installation) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Username,
		&s.Source,
		&s.InstallationID,
		&s.CreatedAt,
		&s.UpdatedAt,
	}
}

// View returns View object for that struct.
func (s *Installation) View() reform.View {
	return InstallationTable
}

// Table returns Table object for that record.
func (s *Installation) Table() reform.Table {
	return InstallationTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Installation) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Installation) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Installation) HasPK() bool {
	return s.ID != InstallationTable.z[InstallationTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Installation) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = int64(i64)
	} else {
		s.ID = pk.(int64)
	}
}

// check interfaces
var (
	_ reform.View   = InstallationTable
	_ reform.Struct = new(Installation)
	_ reform.Table  = InstallationTable
	_ reform.Record = new(Installation)
	_ fmt.Stringer  = new(Installation)
)

func init() {
	parse.AssertUpToDate(&InstallationTable.s, new(Installation))
}
