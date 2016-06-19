package inventory

import "database/sql"

// Databaser allows structs to uniformly interface with a database
type Databaser interface {
	insert(d *sql.DB) (int64, error)
	query(s *sql.Stmt, d *sql.DB) (interface{}, error)
	delete(d *sql.DB) error
	update(d *sql.DB) error
	validateFields(d *sql.DB) error
}
