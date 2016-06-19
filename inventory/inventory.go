package inventory

import (
	"database/sql"
	"fmt"
	"sync"
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
func GetStock(d *sql.DB, u string) (*Stock, error) {
	var err error
	once.Do(func() {
		instance = &Stock{
			DB:  d,
			URL: u,
		}

		// Initialize Tables

		initStmt := `
      create table if not exists categories(
        id integer not null primary key ,
        name text unique
      );
      create table if not exists products(
        id integer not null primary key,
				name text unique not null,
        category_id integer not null,
        price unsigned integer not null,
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
