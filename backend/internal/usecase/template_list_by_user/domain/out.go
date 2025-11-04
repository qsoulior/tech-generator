package domain

type TemplateListByUserOut struct {
	Templates      []Template
	TotalTemplates int64
	TotalPages     int64
}
