package domain

type ProjectListByUserOut struct {
	Owned  []Project
	Shared []Project
}
