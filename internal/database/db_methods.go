package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"tuum.com/internal/models"
	_ "tuum.com/internal/models"
)

type User models.User
type Post models.Post
type Room models.Room
type Comment models.Comment

func CreateUser(username, email, passwordHash string) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)`
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

func CreatePost(title, content string) {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	query := `INSERT INTO posts (user_id, room_id, title, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, title, content, time.Now())
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

func GetUsers() []User {
	db, _ := sql.Open("sqlite3", "./forum.db")
	rows, err := db.Query("SELECT id, username, email, created_at FROM users")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			fmt.Println(err)
			return nil
		}
		users = append(users, user)
	}
	return users
}

func GetUser(id int) User {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	row := db.QueryRow("SELECT id, username, email, created_at FROM users WHERE id = ?", id)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
		fmt.Println(err)
		return User{}
	}
	return user
}

func GetUserByEmail(email string) User {
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	row := db.QueryRow("SELECT id, username, email, created_at FROM users WHERE email = ?", email)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
		fmt.Println(err)
		return User{}
	}
	return user
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