package inventory

import (
	"database/sql"
	"fmt"
)

//Product is the model for a set of items
type Product struct {
	ID       int64
	Price    uint64
	Name     string
	Category string
}

func (p Product) insert(d *sql.DB) (int64, error) {

	q := `INSERT INTO
  products (name, category_id, price) values(?,?,?)
  `
	rows, err := d.Query("SELECT id FROM categories WHERE name='" + p.Category + "'")
	if err != nil {
		return -1, err
	}

	if !rows.Next() {
		return -1, fmt.Errorf("The '%s' category does not exist", p.Category)
	}
	cid := 0
	err = rows.Scan(&cid)
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
