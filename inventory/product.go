package inventory

import (
	"database/sql"
	"fmt"
	"strings"
)

//Product is the model for a set of items
type Product struct {
	ID       int64
	Price    uint64
	Name     string
	Category string
}

func (p Product) validateFields(d *sql.DB) (err error) {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return fmt.Errorf("The name of a product must be provided")
	}

	stmt, err := d.Prepare(
		`SELECT * FROM categories WHERE name=?`
	)

	rows, err := stmt.Query(p.Category)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("The '%s' category does not exist")
	}

	return nil
}

func (p Product) insert(d *sql.DB) (int64, error) {
	q := `INSERT INTO
  product (name, category_id, price) values(?,?,?)
  `
	row := d.QueryRow("SELECT id FROM categories WHERE name=''" + p.Category + "';")

	cid := 0
	err := row.Scan(&cid)
	if err != nil {
		return -1, err
	}

	stmt, err := d.Prepare(q)
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(p.Name, cid, p.Price)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (c Product) query(s *sql.Stmt, d *sql.DB) (interface{}, error) {

	rows, err := s.Query()
	if err != nil {
		return nil, err
	}

	r := []Product{}
	for rows.Next() {
		var id int64
		var n string
		var p uint64
		err = rows.Scan(&id, &n, &p)
		if err != nil {
			return nil, err
		}
		r = append(r, Product{
			ID:    id,
			Name:  n,
			Price: p,
		})
	}

	return r, nil
}

func (c Product) delete(d *sql.DB) error {
	fmt.Println("Delete Called")
	return nil
}

func (c Product) update(d *sql.DB) error {
	fmt.Println("Update Called")
	return nil
}
