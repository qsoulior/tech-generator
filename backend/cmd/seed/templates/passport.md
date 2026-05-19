---
**Сервис:** {{ .service_name }}
**Слаг:** `{{ .service_slug }}`
**Версия:** v{{ .service_version }}
**Владелец:** {{ .service_owner }}
**В эксплуатации с:** {{ .production_date }}
---

# Технический паспорт сервиса «{{ .service_name }}»

> Документ сформирован автоматически {{ date "02.01.2006" now }}.
> Класс нагрузки: **{{ .load_class }}**.

## 1. Назначение и состав

Сервис **{{ .service_name }}** обеспечивает доставку транзакционных и маркетинговых уведомлений по каналам push, SMS и email. Реализован на {{ .runtime }}, развёрнут в кластере из {{ .instances }} реплик.

**Зависимости:**

| Тип | Адрес | Назначение |
| --- | --- | --- |
| База данных | `{{ .dep_db }}` | Персистентное хранилище подписок и журнала отправок |
| Очередь сообщений | `{{ .dep_queue }}` | Асинхронные команды на отправку |
| Кеш | `{{ .dep_cache }}` | Хранение токенов и rate-limit счётчиков |

## 2. Целевые показатели SLO/SLA

| Параметр | Значение |
| --- | --- |
| Доступность (SLA) | **{{ .availability_pct }}** |
| Бюджет ошибок | **{{ .error_budget_min_per_month }} мин / 30 дн** (≈ {{ .error_budget_sci }}) |
| Пропускная способность (цель) | **{{ .rps_target }} RPS** |
| Пиковая нагрузка | **{{ printf "%.0f" .peak_rps }} RPS** |
| Латентность p50 | {{ .latency_p50_ms }} мс |
| Латентность p95 | {{ .latency_p95_ms }} мс |
| Латентность p99 | {{ .latency_p99_ms }} мс |
| Разрыв p99 − p95 | {{ .latency_jump_ms }} мс |

### Формула бюджета ошибок

$$ E_{budget} = (1 - SLA) \times T_{period} = (1 - {{ printf "%.4f" .sla_target }}) \times 43{,}200\,\text{мин} = {{ .error_budget_min_per_month }}\,\text{мин} $$

## 3. Ресурсы

- Количество реплик: **{{ .instances }}**
- Суммарный CPU: **{{ .total_cpu }} vCPU**
- Суммарная память: **{{ .total_memory_gb }} ГБ**

```yaml
service:
  name: {{ .service_slug }}
  version: {{ .service_version }}
  runtime: {{ .runtime }}
  replicas: {{ .instances }}
  resources:
    cpu: {{ .cpu_per_instance }}
    memory: {{ .memory_mb_per_instance }}Mi
  slo:
    availability: {{ printf "%.4f" .sla_target }}
    latency_p99_ms: {{ .latency_p99_ms }}
```

## 4. Чек-лист готовности к продакшену

- [x] Конфигурация для prod утверждена
- [x] Метрики экспортируются в Prometheus
- [x] Алёрты подключены к on-call ротации
- [x] Документация по эксплуатации в Confluence
- [ ] Пройдены нагрузочные испытания на {{ printf "%.0f" .peak_rps }} RPS
- [ ] Согласован план отката
