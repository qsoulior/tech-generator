package main

import (
	"embed"
	"fmt"
)

// templatesFS содержит содержимое всех шаблонов (стандартных и кастомных),
// чтобы тело не приходилось хранить в Go-строках.
//
//go:embed templates/*.md
var templatesFS embed.FS

func readTemplate(name string) string {
	data, err := templatesFS.ReadFile("templates/" + name)
	if err != nil {
		panic(fmt.Sprintf("read embedded template %q: %v", name, err))
	}
	return string(data)
}
