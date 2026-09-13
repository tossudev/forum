package database

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

// isUniqueErr checks if the error violates unique constraint
func IsUniqueErr(err error) bool {
	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return true
	}

	return false
}

// isForeignKeyError checks if the error is a foreign key constraint error.
func IsForeignKeyError(err error) bool {
	var sqliteErr sqlite3.Error

	// Check if err is a foreign key constraint error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
		return true
	}

	return false
}
