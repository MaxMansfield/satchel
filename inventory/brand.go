package inventory

import (
	"database/sql"
	"fmt"
	"strings"
)

// Brand hold data for the brands table rows
type Brand struct {
	ID   int64
	Name string
}

func (b Brand) insert(d *sql.DB) (int64, error) {

	b.Name = strings.TrimSpace(b.Name)
	if b.Name == "" {
		return -1, fmt.Errorf("The name of a brand cannot be blank")
	}

	stmt, err := d.Prepare("INSERT INTO brands (name) values(?)")
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(b.Name)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (b Brand) query(s *sql.Stmt, d *sql.DB) (interface{}, error) {

	rows, err := s.Query()
	if err != nil {
		return nil, err
	}

	r := []Brand{}
	for rows.Next() {
		var id int64
		var n string
		err = rows.Scan(&id, &n)
		if err != nil {
			return nil, err
		}
		r = append(r, Brand{
			ID:   id,
			Name: n,
		})
	}

	return r, nil
}

func (b Brand) delete(d *sql.DB) error {
	fmt.Println("Delete Called")
	return nil
}

func (b Brand) update(d *sql.DB) error {
	fmt.Println("Update Called")
	return nil
}
