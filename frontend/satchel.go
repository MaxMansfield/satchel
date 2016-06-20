package frontend

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/MaxMansfield/satchel/inventory"
	"github.com/urfave/cli"
)

var (
	instance *Satchel
	once     sync.Once
)

// Satchel encapsulates the functions of the CLI
type Satchel struct {
	App   *cli.App
	Stock *inventory.Stock
}

// Captialize strings
func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " a an on the to "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

// Run begins the application
func (s Satchel) Run(args []string) error {

	s.App.Commands = []cli.Command{
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
						n := properTitle(c.String("name"))

						cat := inventory.Category{
							Name: n,
						}

						id, err := s.Stock.Add(cat)
						if err != nil {
							cli.ShowAppHelp(c)
							return err
						}

						fmt.Printf("Category #%03d  - '%s' - has been successfully inserted\n", id, cat.Name)

						return nil
					},
				}, // Category
				{
					Name:    "brand",
					Aliases: []string{"b"},
					Usage:   "add a new brand",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "Specify the name of a brand",
						},
					},
					Action: func(c *cli.Context) error {
						n := properTitle(c.String("name"))
						b := inventory.Brand{
							Name: n,
						}

						id, err := s.Stock.Add(b)
						if err != nil {
							cli.ShowAppHelp(c)
							return err
						}

						fmt.Printf("Brand #%03d  - '%s' - has been successfully inserted\n", id, b.Name)

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
						cli.StringFlag{
							Name:  "brand, b",
							Usage: "The brand of the product",
						},
						cli.Float64Flag{
							Name:  "price, p",
							Usage: "The base price of a product",
						},
					},
					Action: func(c *cli.Context) error {

						n := properTitle(c.String("name"))
						b := properTitle(c.String("brand"))
						cat := properTitle(c.String("category"))

						pro := inventory.Product{
							Name:     n,
							Category: cat,
							Brand:    b,
							Price:    uint64(100 * c.Float64("price")),
						}

						id, err := s.Stock.Add(pro)
						if err != nil {
							return err
						}

						fmt.Printf("Product #%03d  - %s: '%s %s' - has been added\n", id, pro.Category, pro.Brand, pro.Name)

						return nil
					},
				}, // Product
			},
		}, // Add
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List categories, products and items",
			Action:  s.Stock.List,
			Subcommands: []cli.Command{
				{
					Name:    "products",
					Aliases: []string{"p"},
					Usage:   "Print a list of products in the inventory",
					Action:  s.Stock.List,
				}, // Products
				{
					Name:    "brands",
					Aliases: []string{"b"},
					Usage:   "Print a list of brands in the inventory",
					Action:  s.Stock.ListBrands,
				}, // Brands
				{
					Name:    "category",
					Aliases: []string{"c"},
					Usage:   "Print a list of categories in the inventory",
					Action:  s.Stock.List,
				},
			}, // Subcommands
		}, // list
	}

	return s.App.Run(args)
}

// GetCLI attaches commands to a database and app object
func GetCLI(a *cli.App, d *sql.DB) *Satchel {
	once.Do(func() {
		stock, err := inventory.GetStock(d)
		if err != nil {
			return
		}

		instance = &Satchel{
			App:   a,
			Stock: stock,
		}
	})

	return instance
}
