package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"tuum.com/internal/models"
	_ "tuum.com/internal/models"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
}

type User models.User
type Room models.Room
type Post models.Post
type Comment models.Comment

func CreateUser(username, email, passwordHash string) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `INSERT INTO users (username, email, password_hash, status, created_at) VALUES (?, ?, ?, 'user', ?)`
	_, err := db.Exec(query, username, email, passwordHash, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("User created successfully")
}

func CreateRoom(name, description string) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `INSERT INTO rooms (name, description, created_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, name, description, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Room created successfully")
}

func CreatePost(userID, roomID int, title, content string) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `INSERT INTO posts (user_id, room_id, title, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, userID, roomID, title, content, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Post created successfully")
}

func CreateComment(postID, userID int, content string) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, postID, userID, content, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Comment created successfully")
}

func AddPost(post models.Post) error {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `INSERT INTO posts (user_id, room_id, title, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, post.UserID, post.RoomID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Post created successfully")
	return nil
}

func GetUserByEmail(email string) (User, error) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	// Added status to the SELECT query
	row := db.QueryRow("SELECT id, username, email, status, created_at FROM users WHERE email = ?", email)
	var user User
	// Added &user.Status to the row.Scan to capture the status from the query
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Status, &user.CreatedAt); err != nil {
		fmt.Println(err)
		return User{}, err
	}
	return user, nil
}

func GetStatusByEmail(email string) (string, error) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	row := db.QueryRow("SELECT status FROM users WHERE email = ?", email)
	var status string
	if err := row.Scan(&status); err != nil {
		fmt.Println(err)
		return "", err
	}
	return status, nil

}
func GetRooms() []Room {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	rows, err := db.Query("SELECT id, name, description, created_at FROM rooms")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.CreatedAt); err != nil {
			fmt.Println(err)
			return nil
		}
		rooms = append(rooms, room)
	}
	return rooms
}

func GetRoomByName(name string) Room {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	row := db.QueryRow("SELECT id, name, description, created_at FROM rooms WHERE name = ?", name)
	var room Room
	if err := row.Scan(&room.ID, &room.Name, &room.Description, &room.CreatedAt); err != nil {
		fmt.Println(err)
		return Room{}
	}
	return room
}
func GetRoomIdByName(name string) int {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	row := db.QueryRow("SELECT id FROM rooms WHERE name = ?", name)
	var room Room
	err := row.Scan(&room.ID)
	if err != nil {
		fmt.Println(err)
		return room.ID
	}
	return room.ID
}

func GetPosts() []Post {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	rows, err := db.Query("SELECT id, user_id, room_id, title, content, created_at FROM posts")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.RoomID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			fmt.Println(err)
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

func GetComments() []Comment {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			fmt.Println(err)
			return nil
		}
		comments = append(comments, comment)
	}
	return comments
}

func CheckUserExists(username, email string) bool {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	row := db.QueryRow("SELECT username, email FROM users WHERE username = ? OR email = ?", username, email)
	var user User
	err := row.Scan(&user.Username, &user.Email)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func CheckRoomExists(name string) bool {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	row := db.QueryRow("SELECT name FROM rooms WHERE name = ?", name)
	var room Room
	err := row.Scan(&room.Name)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

type DatabaseContent struct {
	Rooms []Room
	// Posts    []Post
	// Comments []Comment
}

func GetDatabaseForTuum() DatabaseContent {
	rooms := GetRooms()
	if rooms == nil {
		fmt.Println("rooms error")
		return DatabaseContent{}
	}

	// posts := GetPosts()
	// if posts == nil {
	// 	fmt.Println("posts error")
	// 	return DatabaseContent{}
	// }

	// comments := GetComments()
	// if comments == nil {
	// 	fmt.Println("comments error")
	// 	return DatabaseContent{}
	// }
	fmt.Println("Database content fetched successfully")

	return DatabaseContent{
		Rooms: rooms,
		// Posts:    posts,
		// Comments: comments,
	}

}

func DeleteAccountByEmail(email string) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `DELETE FROM users WHERE email = ?`
	result, err := db.Exec(query, email)
	if err != nil {
		log.Printf("Error executing delete query for email %s: %v", email, err)
		return fmt.Errorf("error executing delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected for email %s: %v", email, err)
		return fmt.Errorf("error fetching rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("No user found with the given email: %s", email)
		return errors.New("no user found with the given email")
	}

	log.Printf("Account deleted successfully for email: %s", email)
	return nil
}
