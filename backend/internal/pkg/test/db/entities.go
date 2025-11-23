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

type Project struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	AuthorID int64  `db:"author_id"`
}

type ProjectUser struct {
	ProjectID int64  `db:"project_id"`
	UserID    int64  `db:"user_id"`
	Role      string `db:"role" fake:"{randomstring:[read,write,maintain]}"`
}

type Template struct {
	ID            int64      `db:"id"`
	Name          string     `db:"name"`
	IsDefault     bool       `db:"is_default"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	ProjectID     *int64     `db:"project_id"`
	AuthorID      *int64     `db:"author_id"`
	LastVersionID *int64     `db:"last_version_id"`
}

type TemplateUser struct {
	TemplateID int64  `db:"template_id"`
	UserID     int64  `db:"user_id"`
	Role       string `db:"role" fake:"{randomstring:[read,write]}"`
}

type Version struct {
	ID         int64     `db:"id"`
	Number     int64     `db:"number"`
	TemplateID int64     `db:"template_id"`
	AuthorID   *int64    `db:"author_id"`
	CreatedAt  time.Time `db:"created_at"`
	Data       []byte    `db:"data"`
}

type Variable struct {
	ID         int64   `db:"id"`
	VersionID  int64   `db:"version_id"`
	Name       string  `db:"name"`
	Type       string  `db:"type" fake:"{randomstring:[integer,float,string]}"`
	Expression *string `db:"expression"`
	IsInput    bool    `db:"is_input"`
}

type Constraint struct {
	ID         int64  `db:"id"`
	VariableID int64  `db:"variable_id"`
	Name       string `db:"name"`
	Expression string `db:"expression"`
	IsActive   bool   `db:"is_active"`
}

type Task struct {
	ID        int64      `db:"id"`
	VersionID int64      `db:"version_id"`
	Status    string     `db:"status" fake:"{randomstring:[created,in_progress,succeed,failed]}"`
	Payload   []byte     `db:"payload"`
	ResultID  *int64     `db:"result_id"`
	Error     []byte     `db:"error"`
	CreatorID int64      `db:"creator_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type Result struct {
	ID   int64  `db:"id"`
	Data []byte `db:"data"`
}
