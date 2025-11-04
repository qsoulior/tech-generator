package domain

type ProjectListByUserOut struct {
	Projects      []Project
	TotalProjects int64
	TotalPages    int64
}
