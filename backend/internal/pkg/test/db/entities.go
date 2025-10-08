package test_db

import "time"

type User struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  []byte     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type Folder struct {
	ID           int64  `db:"id"`
	ParentID     *int64 `db:"parent_id"`
	Name         string `db:"name"`
	AuthorID     int64  `db:"author_id"`
	RootAuthorID int64  `db:"root_author_id"`
}

type FolderUser struct {
	FolderID int64  `db:"folder_id"`
	UserID   int64  `db:"user_id"`
	Role     string `db:"role"`
}
