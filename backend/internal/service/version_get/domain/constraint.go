package domain

type Constraint struct {
	ID         int64
	VariableID int64
	Name       string
	Expression string
	IsActive   bool
}
