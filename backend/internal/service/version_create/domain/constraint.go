package domain

type Constraint struct {
	Name       string
	Expression string
	IsActive   bool
}

type ConstraintToCreate struct {
	VariableID int64
	Name       string
	Expression string
	IsActive   bool
}
