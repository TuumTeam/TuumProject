package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"tuum.com/internal/models"
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	createTables(db)
}

func createTables(db *sql.DB) {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        status TEXT DEFAULT 'user' NOT NULL,
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

func Login(email string, password string) (bool, error) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `SELECT * FROM users WHERE email = ?`
	row := db.QueryRow(query, email)
	user := models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with the provided email
			return false, errors.New("no user found with the provided email")
		}
		// Other error
		return false, err
	}
	if password == user.Password {
		return true, nil
	} else {
		// Incorrect password
		return false, errors.New("incorrect password")
	}
}
