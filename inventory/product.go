package inventory

import (
	"database/sql"
	"fmt"
	"sort"
)

// Products is a sortable slice of Product
type Products []Product

//Len returns the length of Products
func (p Products) Len() int {
	return len(p)
}

//Less determines if one index is less than another
func (p Products) Less(i, j int) bool {
	strs := []string{
		p[i].Brand,
		p[j].Brand,
	}
	sort.Strings(strs)
	return strs[0] == p[i].Brand
}

//Swap will swap two products
func (p Products) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

//Product is the model for a set of items
type Product struct {
	ID       int64
	Price    uint64
	Name     string
	Category string
	Brand    string
}

func (p Product) insert(d *sql.DB) (int64, error) {

	q := `INSERT INTO
  products (name, price, category_id, brand_id) values(?,?,?,?)
  `
	tx, err := d.Begin()
	if err != nil {
		return -1, err
	}

	stmt, err := tx.Prepare("SELECT id FROM categories WHERE name = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(p.Category)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if !rows.Next() {
		return -1, fmt.Errorf("The '%s' category does not exist", p.Category)
	}
	if err = rows.Err(); err != nil {
		return -1, err
	}

	cid := 0
	err = rows.Scan(&cid)
	if err != nil {
		return -1, err
	}

	stmt, err = tx.Prepare("SELECT id FROM brands WHERE name = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rows, err = stmt.Query(p.Brand)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if !rows.Next() {
		return -1, fmt.Errorf("The '%s' brand does not exist", p.Brand)
	}
	if err = rows.Err(); err != nil {
		return -1, err
	}

	bid := 0
	err = rows.Scan(&bid)
	if err != nil {
		return -1, err
	}

	stmt, err = tx.Prepare(q)
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(p.Name, p.Price, cid, bid)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	tx.Commit()
	return id, nil
}

func (p Product) query(s *sql.Stmt, d *sql.DB) (interface{}, error) {

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

func (p Product) delete(d *sql.DB) error {
	fmt.Println("Delete Called")
	return nil
}

func (p Product) update(d *sql.DB) error {
	fmt.Println("Update Called")
	return nil
}
