package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "akash"
	password = "25651640"
	dbname   = "gophercises_phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	// db, err := sql.Open("postgres", psqlInfo)
	// must(err)

	// err = resetDB(db, dbname)
	// must(err)
	// db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Ping())
	must(createTable(db))
	// id, err := insertPhone(db, "(123)456-7892")
	// must(err)
	// fmt.Println("id=", id)

	number, err := getPhone(db, 4)
	must(err)
	fmt.Println("Number is ...", number)

	phones, err := allPhones(db)
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on .. %+v\n", p)
		number := normalize(p.number)
		if number != p.number {
			fmt.Println("Updating or removing...", number)
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {

			}else{
				
			}
		} else {
			fmt.Println("No changes requried")
		}
	}

}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE id=$1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

type phone struct {
	id     int
	number string
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&p.id, &p.number)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func allPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func createTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255) 
		)`
	_, err := db.Exec(statement)
	return err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func normalize(phone string) string {

	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(phone, "")
}

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, ch := range phone {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}
// 	return buf.String()
// }
