package repositories

import (
	"database/sql"
)

type PostRepository struct {
	DB *sql.DB
}
