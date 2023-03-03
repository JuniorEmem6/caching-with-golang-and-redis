package data

import (
	"context"
	"database/sql" // add this
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Product struct {
	Name       string `json:"name"`
	Price      string `json:"price"`
	Descrition string `json:"description"`
}

var db *sql.DB
var err error

func Connect() {
	var constr string = "postgres://postgres:Patience10();@localhost/cache?sslmode=disable"
	db, err = sql.Open("postgres", constr)

	if err != nil {
		fmt.Println(err.Error())
	}

	if err1 := db.Ping(); err1 != nil {
		log.Fatalf("unable to reach database: %v", err1)
	}

	fmt.Println("Mysql Database Connected")
}

func SelectProduct(id int) (string, string, string, error) {
	var query string = "SELECT name, price, description FROM products where id = $1"
	var name string
	var price string
	var description string

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return name, price, description, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	if err := row.Scan(&name, &price, &description); err != nil {
		return name, price, description, err
	}

	return name, price, description, nil
}

func InsertProduct(product *Product) error {
	var query string = "INSERT INTO products (name, price, description) VALUES ($1, $2, $3) RETURNING *"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, product.Name, product.Price, product.Descrition)
	if err != nil {
		return errors.New("could not insert row")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not get affected rows")
	}

	// we can log how many rows were inserted
	fmt.Println("inserted", rowsAffected, "rows")
	return nil
}
