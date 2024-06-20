package database

import (
	"database/sql"
	"fmt"
	"log"
	"tuum.com/internal/models"
	"tuum.com/internal/services"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./database/users.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func createTable() {
	query := `
 CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL
 );`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
}

func AddUser(user models.User) error {
	user.Password, _ = services.HashPassword(user.Password)

	sqlStatement := `
 INSERT INTO users (username, email, password)
 VALUES (?, ?, ?)`
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return err
	}
	fmt.Println("New record ID is:", id)
	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	row := db.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Login(email, password string) (bool, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return false, err
	}

	return services.CheckPasswordHash(password, user.Password), nil
}
