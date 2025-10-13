package domain

type Template struct {
	Name         string
	IsDefault    bool
	FolderID     int64
	AuthorID     int64
	RootAuthorID int64
}
