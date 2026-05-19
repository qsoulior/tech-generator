# Релиз {{ .release_version }}

**Дата:** {{ .release_date }} ⋅ **Релиз-мастер:** {{ .release_owner }}

**Breaking changes:** {{ .has_breaking }}

## Что вошло

| Категория | Количество |
| --- | --- |
| Новые фичи | {{ .features_added }} |
| Исправленные баги | {{ .bugs_fixed }} |
| **Всего изменений** | **{{ .total_changes }}** |

## Идентификатор для CI/CD

`{{ .release_slug }}`
