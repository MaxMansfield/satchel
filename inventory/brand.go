package inventory

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
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

//ListBrands list all brands to stdout in a table
// FOR inventory.Stock!!! Read carefully
func (s Stock) ListBrands(c *cli.Context) error {
	t, err := s.DB.Begin()
	if err != nil {
		return err
	}

	r, err := t.Query(`SELECT id,name FROM brands ORDER BY name ASC`)
	if err != nil {
		return err
	}
	defer r.Close()

	pros := make(map[string]Products)
	for r.Next() {
		var id int64
		var n string
		err = r.Scan(&id, &n)
		if err != nil {
			return err
		}

		rows, err := t.Query(fmt.Sprintf(
			`SELECT name, price FROM products WHERE brand_id = %d ORDER BY name DESC;`,
			id,
		))
		if err != nil {
			return err
		}
		defer r.Close()

		for rows.Next() {
			var pn string
			var pp uint64

			err = rows.Scan(&pn, &pp)
			if err != nil {
				return err
			}

			pro := Product{
				Name:  pn,
				Price: pp,
			}
			pros[n] = append(pros[n], pro)
		}

		if err = rows.Err(); err != nil {
			return err
		}

	}
	if err = r.Err(); err != nil {
		return err
	}

	t.Commit()

	var keys []string
	for k := range pros {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	table := tablewriter.NewWriter(os.Stdout)

	prolen := 0
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
				v.Name,
				fmt.Sprintf("$%0.2f", float64(v.Price)/100),
			})
		}
	}

	table.SetHeader([]string{"Brand", "Model", "Price"})
	table.SetFooter([]string{
		"",
		"Brands",
		fmt.Sprintf("%02d", prolen),
	})
	table.SetBorder(false) // Set Border to false

	table.Render()

	return nil
}
