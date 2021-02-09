package glowlib

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// ===========================================================================================================================================

// NullUnixTime represents a Unix time value that may be null. NullUnixTime implements the
// sql.Scanner interface so it can be used as a scan destination, similar to sql.NullString
type NullUnixTime struct {
	UnixTime float64
	Valid    bool // Valid is true if Time is not NULL
}

// Scan converts from interface
func (nt *NullUnixTime) Scan(value interface{}) error {
	var f64 float64
	f64, nt.Valid = value.(float64) // EXTRACT (EPOCH FROM x) returns a double-precision value
	nt.UnixTime = f64               // we keep it as a double to retain the sub-second accuracy for chunking queries
	return nil
}

// Value converts from driver value
func (nt NullUnixTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.UnixTime, nil
}

// ===========================================================================================================================================

// NullableU64 represents a u64 that can also be null
type NullableU64 struct {
	ValueU64 uint64
	Valid    bool // Valid is true if Time is not NULL
}

// Scan converts from interface
func (nt *NullableU64) Scan(value interface{}) error {
	var vu64 uint64
	vu64, nt.Valid = value.(uint64)
	nt.ValueU64 = vu64
	return nil
}

// Value converts from driver value
func (nt NullableU64) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.ValueU64, nil
}

// ===========================================================================================================================================

/*
PrecompileSQL reflects across a structure, finds sql.Stmt pointers with a sql tag and Prepare()s the sql, storing
the precompiled statement result back in the passed block

type sqlPrecomp struct {
	GrabWhereXandY     *sql.Stmt `sql:"SELECT * FROM tbl.football WHERE x=$1 and y=$2"`
}

var mySql sqlPrecomp

PrecompileSQL(&mySql)
*/
func precompileSQLInternal(db *sql.DB, statementField interface{}) error {

	smType := reflect.TypeOf(statementField).Elem()
	smValue := reflect.ValueOf(statementField).Elem()

	for i := 0; i < smType.NumField(); i++ {
		fieldType := smType.Field(i)

		field := smValue.Field(i)

		sqlStatement := fieldType.Tag.Get("sql")
		itype := reflect.TypeOf(field.Interface()).String()

		if sqlStatement != "" && itype == "*sql.Stmt" {

			stmt, err := db.Prepare(sqlStatement)
			if err != nil {
				return err
			}

			field.Set(reflect.ValueOf(stmt))
		}
	}

	return nil
}

// PrecompileSQLEx allows you to pass in the db connection you want to use
func PrecompileSQL(db *sql.DB, statementField interface{}) error {
	return precompileSQLInternal(db, statementField)
}

// Individual sql.Stmt compilation
func PrecompileOne(db *sql.DB, sqlstr string, errLog *zap.Logger) (*sql.Stmt, error) {

	stmt, err := db.Prepare(sqlstr)
	if err != nil {
		if errLog != nil {
			errLog.Panic("PrecompileOne failed",
				zap.Error(err))
		}

		return nil, err
	}

	return stmt, nil
}

// yay for flatbuffers encoding booleans as bytes for no good reason
func ByteFromBool(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func StringFromUintArray(a []uint64) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.FormatUint(v, 10)
	}
	return strings.Join(b, ",")
}
