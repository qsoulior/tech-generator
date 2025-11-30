package domain

type TaskCreateIn struct {
	VersionID int64
	CreatorID int64
	Payload   map[string]string
}
