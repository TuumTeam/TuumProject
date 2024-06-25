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

import (
	"database/sql"
	"fmt"
	"tuum.com/internal/models"
	"tuum.com/internal/services"
)

func init() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	createTables(db)
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

func Login(email string, password string) bool {
	db, _ := sql.Open("sqlite3", "./forum.db")
	query := `SELECT * FROM users WHERE email = ?`
	row := db.QueryRow(query, email)
	user := models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return false
	}
	if services.CheckPasswordHash(password, user.Password) {
		return true
	} else {
		return false
	}
}
