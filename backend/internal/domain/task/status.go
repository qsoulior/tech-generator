package task_domain

type Status string

const (
	StatusCreated    Status = "created"
	StatusInProgress Status = "in_progress"
	StatusSucceed    Status = "succeed"
	StatusFailed     Status = "failed"
)

var statusSet = map[Status]struct{}{
	StatusCreated:    {},
	StatusInProgress: {},
	StatusSucceed:    {},
	StatusFailed:     {},
}

func (s Status) Valid() bool {
	_, found := statusSet[s]
	return found
}
