package database

import (
	"database/sql"
	"fmt"
	"log"
	"tuum.com/internal/models"
	"tuum.com/internal/services"
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
	createTables(db)
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

func Login(email, password string) (bool, error) {
	user, _ := GetUserByEmail(email)
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	return services.CheckPasswordHash(password, user.Password), nil
}
func createTables(db *sql.DB) {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	roomTable := `
    CREATE TABLE IF NOT EXISTS rooms (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	postTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        room_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(user_id) REFERENCES users(id),
        FOREIGN KEY(room_id) REFERENCES rooms(id)
    );`

	commentTable := `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(post_id) REFERENCES posts(id),
        FOREIGN KEY(user_id) REFERENCES users(id)
    );`

	_, err := db.Exec(userTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(roomTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(postTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(commentTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Tables created successfully")
}
