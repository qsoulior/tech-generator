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
	Role     string `db:"role" fake:"{randomstring:[read,write,maintain]}"`
}

type Template struct {
	ID            int64      `db:"id"`
	Name          string     `db:"name"`
	IsDefault     bool       `db:"is_default"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	FolderID      *int64     `db:"folder_id"`
	AuthorID      *int64     `db:"author_id"`
	RootAuthorID  *int64     `db:"root_author_id"`
	LastVersionID *int64     `db:"last_version_id"`
}

type TemplateUser struct {
	TemplateID int64  `db:"template_id"`
	UserID     int64  `db:"user_id"`
	Role       string `db:"role" fake:"{randomstring:[read,write]}"`
}

type TemplateVersion struct {
	ID         int64     `db:"id"`
	Number     int64     `db:"number"`
	TemplateID int64     `db:"template_id"`
	AuthorID   *int64    `db:"author_id"`
	CreatedAt  time.Time `db:"created_at"`
	Data       []byte    `db:"data"`
}

type Variable struct {
	ID         int64  `db:"id"`
	VersionID  int64  `db:"version_id"`
	Name       string `db:"name"`
	Type       string `db:"type" fake:"{randomstring:[integer,float,string]}"`
	Expression string `db:"expression"`
}

type VariableConstraint struct {
	ID         int64  `db:"id"`
	VariableID int64  `db:"variable_id"`
	Name       string `db:"name"`
	Expression string `db:"expression"`
	IsActive   bool   `db:"is_active"`
}
