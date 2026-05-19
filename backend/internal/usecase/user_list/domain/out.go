package domain

type UserListOut struct {
	Users      []User
	TotalUsers int64
	TotalPages int64
}
