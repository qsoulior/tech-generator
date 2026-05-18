package domain

type Template struct {
	Name      string
	IsDefault bool
	ProjectID int64
	AuthorID  int64
}

type SourceTemplate struct {
	ID            int64
	IsDefault     bool
	LastVersionID *int64
}
