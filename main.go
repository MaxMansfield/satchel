package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MaxMansfield/satchel/inventory"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

//Compile Time Constants
const (
	// 	BUILD_NAME- The name of the application to be passed in at build
	BUILD_NAME = ""
	// 	BUILD_VERSION- The version of the build - also passed in at build
	BUILD_VERSION = ""
	// 	BUILD_TYPE - The type of build such as release, test or debug
	BUILD_TYPE = ""
	// BUILD_TIME - The time that this source was built
	BUILD_TIME = ""
)

//Defaults
const (
	// The default sqlite db location
	D_DBFile = "db/inventory.db"
)

func main() {
	defer fmt.Println("Exiting...")

	var satchel *inventory.Stock

	u := D_DBFile

	fmt.Printf("Opening Database file '%s'...\n", u)
	db, err := sql.Open("sqlite3", u)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	satchel, err = inventory.GetStock(db, u)
	if err != nil {
		panic(err)
	}

	a := cli.NewApp()
	a.Name = BUILD_NAME
	a.Version = BUILD_VERSION
	a.Author = "Max Mansfield"
	a.Email = "max.m.mansfield@gmail.com"
	a.Usage = "A CLI tool to inventory sellable items"
	a.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a category, item or type to the inventory",
			Subcommands: []cli.Command{
				{
					Name:    "category",
					Aliases: []string{"c"},
					Usage:   "add a new category",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "Specify the name of a category",
						},
					},
					Action: func(c *cli.Context) error {
						if c.NumFlags() < 1 {
							cli.ShowAppHelp(c)
							return errors.New("A name must be supplied to add a category")
						}

						n := c.String("name")
						cat := inventory.Category{
							Name: n,
						}

						id, err := satchel.Add(cat)
						if err != nil {
							return err
						}

						fmt.Printf("Category #%03d  - '%s' - has been successfully inserted\n", id, cat.Name)

						return nil
					},
				}, // Category
				{
					Name:    "product",
					Aliases: []string{"p"},
					Usage:   "add a new product",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "The name of the product",
						},
						cli.StringFlag{
							Name:  "category, c",
							Usage: "The category that a product belongs to",
						},
						cli.Float64Flag{
							Name:  "price, p",
							Usage: "The base price of a product",
						},
					},
					Action: func(c *cli.Context) error {
						fmt.Println(c.NumFlags())

						pro := inventory.Product{
							Name:     c.String("name"),
							Category: c.String("category"),
							Price:    uint64(100 * c.Float64("price")),
						}

						id, err := satchel.Add(pro)
						if err != nil {
							return err
						}

						fmt.Printf("Product #%03d  - '%s' - has been successfully inserted\n", id, pro.Name)

						return nil
					},
				}, // Product
			},
		}, // Add
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List categories, products and items",
			Subcommands: []cli.Command{
				{
					Name:    "categories",
					Aliases: []string{"c"},
					Usage:   "list categories",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "Select a category by name",
						},
						cli.StringFlag{
							Name:  "query, q",
							Value: "",
							Usage: "Use a direct sql query",
						},
					},
					Action: func(c *cli.Context) error {
						n := strings.TrimSpace(c.String("name"))
						sql := "SELECT * FROM categories "
						cats := []inventory.Category{}
						wanted := inventory.Category{
							Name: n,
						}

						query := func(str string) error {
							stmt, err := db.Prepare(str)
							if err != nil {
								return err
							}
							dbsrs, err := satchel.Query(stmt, wanted)
							if err != nil {
								return err
							}

							cats = append(cats, dbsrs.([]inventory.Category)...)

							return nil
						}

						if n != "" {
							sql += "WHERE name = " + n
						}

						q := c.String("query")
						if q != "" {
							sql = q
						}

						query(sql + " ORDER BY name ASC;")

						fmt.Printf("Categories (%02d):\n", len(cats))
						for i, v := range cats {
							fmt.Printf("\t#%d %s\n", i+1, v.Name)
						}

						return nil
					},
				},
			}, // categories
		}, // list
	}

	if err := a.Run(os.Args); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}

	//
	// stmt, err = db.Prepare("select name from foo where id = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// var name string
	// err = stmt.QueryRow("3").Scan(&name)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(name)
	//
	// _, err = db.Exec("delete from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// _, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// rows, err = db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

//
// // AddCategory Add a Category to the database
// func AddCategory(n string) error {
//   tx, err :=
// 	s, err = tx.Prepare("insert into foo(id, name) values(?, ?)")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
//
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
// 	if err != nil {
// 		return err
// 	}
//
// 	tx.Commit()
// }
//
// func ListCatagories() error {
// 	rows, err := db.Query("select * from catagories")
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()
//
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		err = rows.Scan(&id, &name)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println(id, name)
// 	}
//
// 	err = rows.Err()
// 	if err != nil {
// 		return err
// 	}
// }
