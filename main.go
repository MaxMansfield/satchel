package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/MaxMansfield/satchel/frontend"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

//Compile Time Constants
const (
	// 	BUILD_NAME- The name of the application to be passed in at build
	BUILDNAME = ""
	// 	BUILD_VERSION- The version of the build - also passed in at build
	BUILDVERSION = ""
	// 	BUILD_TYPE - The type of build such as release, test or debug
	BUILDTYPE = ""
	// BUILD_TIME - The time that this source was built
	BUILDTIME = ""
)

//Defaults
const (
	// The default sqlite db location
	DDBFile = "db/inventory.db"
)

func main() {
	u := DDBFile
	db, err := sql.Open("sqlite3", u)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	a := cli.NewApp()
	a.Name = BUILDNAME
	a.Version = BUILDVERSION
	a.Author = "Max Mansfield"
	a.Email = "max.m.mansfield@gmail.com"
	a.Usage = "A CLI tool to inventory sellable items"

	satchel := frontend.GetCLI(a, db)
	if satchel == nil {
		panic(fmt.Errorf("Unable to allocate CLI"))
	}

	if err = satchel.Run(os.Args); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}
