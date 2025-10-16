package variable_domain

type Type string

const (
	TypeString  Type = "string"
	TypeInteger Type = "integer"
	TypeFloat   Type = "float"
)

var typeSet = map[Type]struct{}{
	TypeString:  {},
	TypeInteger: {},
	TypeFloat:   {},
}

func (r Type) Valid() bool {
	_, found := typeSet[r]
	return found
}
