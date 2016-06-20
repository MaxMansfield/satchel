package inventory

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var (
	instance *Stock
	once     sync.Once
)

// Stock manages the Databaser items and the database itself
type Stock struct {
	DB  *sql.DB
	URL string
}

// GetStock returns an instance of the Stock singleton
// u - the url of the sqlite database
// d - a pointer to the sql.DB object to user
func GetStock(d *sql.DB) (*Stock, error) {
	var err error
	once.Do(func() {
		instance = &Stock{
			DB: d,
		}
		// Initialize Tables

		initStmt := `
      create table if not exists categories(
        id integer not null primary key ,
        name text unique
      );
			create table if not exists brands (
				id integer not null primary key,
				name text unique not null
			);
      create table if not exists products(
        id integer not null primary key,
				name text unique not null,
        price unsigned integer not null,
				brand_id integer not null,
        category_id integer not null,
				foreign key(brand_id) references brands(id)
				foreign key(category_id) references categories(id)
      );
      create table if not exists stock(
        id integer not null primary key,
        product_id integer not null,
        foreign key(product_id) references products(id)
      );
	  `
		_, err = d.Exec(initStmt)
		if err != nil {
			err = fmt.Errorf("%q: %s\n", err, initStmt)
		}

	})

	return instance, err
}

// Add calls a Databasers insert function and passes the DB of the Stock
func (s Stock) Add(d Databaser) (int64, error) {
	return d.insert(s.DB)
}

// Edit calls a Databasers update function and passes the DB of the stock
func (s Stock) Edit(d Databaser) error {
	return d.update(s.DB)
}

//Remove calls a Databasers delete function and passes the DB of the stock
func (s Stock) Remove(d Databaser) error {
	return d.delete(s.DB)
}

// Query calls Databasers query function and passes the DB of the stock
// q - The query to
func (s Stock) Query(q *sql.Stmt, d Databaser) (interface{}, error) {
	return d.query(q, s.DB)
}

// List DB to stdout in a table
func (s Stock) List(c *cli.Context) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query(`
              SELECT
              id,
              name,
              price,
              category_id,
              brand_id
              FROM products;
            `)
	if err != nil {
		return err
	}
	defer rows.Close()

	pros := make(map[string]Products)
	for rows.Next() {
		var id int64
		var n string
		var p uint64
		var cid int64
		var bid int64

		err = rows.Scan(&id, &n, &p, &cid, &bid)
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(
			`SELECT name FROM categories WHERE id = ? ORDER BY name ASC`,
		)
		if err != nil {
			return err
		}
		defer stmt.Close()

		row := stmt.QueryRow(cid)
		if row == nil {
			panic(fmt.Errorf("'%s' Product Found without category", n))
		}

		var cat string
		row.Scan(&cat)

		stmt, err = tx.Prepare(
			`SELECT name FROM brands WHERE id = ? ORDER BY name ASC;`,
		)
		if err != nil {
			return err
		}
		defer stmt.Close()

		row = stmt.QueryRow(bid)
		if row == nil {
			panic(fmt.Errorf("'%s' Product Found without brand", n))
		}

		var bra string
		row.Scan(&bra)

		pro := Product{
			ID:       id,
			Name:     n,
			Price:    p,
			Category: cat,
			Brand:    bra,
		}

		pros[cat] = append(pros[cat], pro)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	tx.Commit()

	var keys []string
	for k := range pros {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	prolen := 0
	table := tablewriter.NewWriter(os.Stdout)

	for _, k := range keys {
		sort.Sort(pros[k])
		table.Append([]string{
			k + fmt.Sprintf(" (%02d)", len(pros[k])),
			"",
			"",
		})
		for _, v := range pros[k] {
			prolen++
			table.Append([]string{
				"",
				v.Brand,
				v.Name,
				fmt.Sprintf("$%0.2f", float64(v.Price)/100),
			})
		}
	}

	table.SetHeader([]string{"Category", "Brand", "Model", "Price"})
	table.SetFooter([]string{
		"Categories",
		fmt.Sprintf("%02d", len(pros)+1),
		"Products",
		fmt.Sprintf("%02d", prolen),
	})
	table.SetBorder(false) // Set Border to false

	table.Render()

	return nil
}
