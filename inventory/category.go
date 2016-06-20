package inventory

import (
	"database/sql"
	"fmt"
	"strings"
)

// Category is the top hierarchichal element of the inventory system.
type Category struct {
	ID   int64
	Name string
}

func (c Category) insert(d *sql.DB) (int64, error) {

	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return -1, fmt.Errorf("The name of a category cannot be blank")
	}

	stmt, err := d.Prepare("INSERT INTO categories (name) values(?)")
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(c.Name)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (c Category) query(s *sql.Stmt, d *sql.DB) (interface{}, error) {

	rows, err := s.Query()
	if err != nil {
		return nil, err
	}

	r := []Category{}
	for rows.Next() {
		var id int64
		var n string
		err = rows.Scan(&id, &n)
		if err != nil {
			return nil, err
		}
		r = append(r, Category{
			ID:   id,
			Name: n,
		})
	}

	return r, nil
}

func (c Category) delete(d *sql.DB) error {
	fmt.Println("Delete Called")
	return nil
}

func (c Category) update(d *sql.DB) error {
	fmt.Println("Update Called")
	return nil
}
