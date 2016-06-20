/** list.go contains all the functions that the frontend satchel object needs
 * to list its data to stdout
 */

package frontend

import (
	"fmt"
	"os"
	"sort"

	"github.com/MaxMansfield/satchel/inventory"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

//listProducts list all products to stdout in a table
func (s Satchel) list(c *cli.Context) error {
	tx, err := s.Stock.DB.Begin()
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

	pros := make(map[string]inventory.Products)
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

		pro := inventory.Product{
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

//listBrand list all brands to stdout in a table
func (s Satchel) listBrands(c *cli.Context) error {
	t, err := s.Stock.DB.Begin()
	if err != nil {
		return err
	}

	r, err := t.Query(`SELECT id,name FROM brands ORDER BY name ASC`)
	if err != nil {
		return err
	}
	defer r.Close()

	pros := make(map[string]inventory.Products)
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

			pro := inventory.Product{
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
