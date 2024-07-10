package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	println("LOL NO MAN")

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// get the version

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("ECCOTE LOL ", version)

	// some raw CRUD

	sts := `
	DROP TABLE IF EXISTS cars;
	CREATE TABLE cars(id INTEGER PRIMARY KEY, name TEXT, price INT);
	INSERT INTO cars(name, price) VALUES('Audi',52642);
	INSERT INTO cars(name, price) VALUES('Mercedes',57127);
	INSERT INTO cars(name, price) VALUES('Skoda',9000);
	INSERT INTO cars(name, price) VALUES('Volvo',29000);
	INSERT INTO cars(name, price) VALUES('Bentley',350000);
	INSERT INTO cars(name, price) VALUES('Citroen',21000);
	INSERT INTO cars(name, price) VALUES('Hummer',41400);
	INSERT INTO cars(name, price) VALUES('Volkswagen',21600);
	`
	_, err = db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("table cars created")

	// raw select

	rows, err := db.Query("SELECT * FROM cars")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var name string
		var price int

		err = rows.Scan(&id, &name, &price)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d %s %d\n", id, name, price)
	}

	// prepared statement

	stm, err := db.Prepare("SELECT * FROM cars WHERE id = ?")

	if err != nil {
		log.Fatal(err)
	}

	defer stm.Close()

	var id int
	var name string
	var price int

	cid := 3

	err = stm.QueryRow(cid).Scan(&id, &name, &price)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d %s %d\n", id, name, price)

	// prepared statement, one shot

	row := db.QueryRow("SELECT * FROM cars WHERE id = ?", cid)
	err = row.Scan(&id, &name, &price)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d %s %d\n", id, name, price)

}
