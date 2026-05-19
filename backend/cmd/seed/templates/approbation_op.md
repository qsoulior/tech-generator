# Описание программы «{{ .program_name }}»

| Реквизит | Значение |
| --- | --- |
| Наименование программы | {{ .program_name }} |
| Условное обозначение | {{ .program_designation }} |
| Версия | {{ .program_version }} |
| Дата документа | {{ .document_date }} |
| Шифр документа | {{ .document_code }} |
| Среда разработки | {{ .runtime }} |
| Целевая ОС | {{ .target_os }} |

> Подготовлено в соответствии с **ГОСТ 19.402-78**.

## 1. Общие сведения

Программа «{{ .program_name }}» (условное обозначение `{{ .program_designation }}`) реализована на {{ .runtime }} и предназначена для развёртывания в среде {{ .target_os }} в контейнерах OCI. Сборка осуществляется единым CI-конвейером, артефакт — мульти-архитектурный контейнерный образ `{{ .container_image }}`.

| Параметр | Значение |
| --- | --- |
| Язык программирования | {{ .language }} |
| Среда выполнения | {{ .runtime }} |
| Способ распространения | OCI-образ + Helm-чарт |
| Метод хранения исходного кода | Git, репозиторий `{{ .source_repo }}` |

## 2. Функциональное назначение

Программа автоматизирует доставку транзакционных и маркетинговых уведомлений по каналам push, SMS и email. Основные функции:

- регистрация и хранение подписок получателей;
- маршрутизация события доставки в нужный канал;
- расчёт политик повторных попыток и подавления спама;
- журналирование и выдача отчётности по доставкам.

Программа взаимодействует с внешними системами по протоколам HTTP/2, AMQP 0.9.1 и SMTP.

## 3. Описание логической структуры

Программа состоит из следующих модулей:

| Модуль | Назначение |
| --- | --- |
| `api-gateway` | приём запросов, аутентификация, маршрутизация |
| `policy-engine` | расчёт правил доставки и подавления |
| `dispatcher` | формирование команд на отправку, повторные попытки |
| `adapter-push` | интеграция с APNs/FCM |
| `adapter-sms` | интеграция с SMS-агрегатором |
| `adapter-email` | интеграция с почтовым релеем |
| `admin-ui` | административная панель оператора |

Взаимодействие модулей — асинхронное, через шину RabbitMQ. Топология обменов:

```
api-gateway ── exchange:notify.in ──► policy-engine
policy-engine ── exchange:notify.dispatch ──► dispatcher
dispatcher    ── routing.push ──► adapter-push
dispatcher    ── routing.sms  ──► adapter-sms
dispatcher    ── routing.email──► adapter-email
```

## 4. Используемые технические средства

Минимальные требования к узлу кластера:

| Параметр | Значение |
| --- | --- |
| Архитектура | {{ .cpu_arch }} |
| Количество vCPU на узел | ≥ {{ .min_cpu_per_node }} |
| Оперативная память на узел | ≥ {{ .min_memory_gb_per_node }} ГБ |
| Свободное место на диске | ≥ {{ .min_disk_gb_per_node }} ГБ |
| Сеть | ≥ {{ .min_network_gbps }} Гбит/с, RTT ≤ 5 мс |

Программа развёртывается в кластере Kubernetes версии {{ .k8s_version }} или выше; для управления нагрузкой используется HPA по метрикам RPS и latency.

## 5. Вызов и загрузка

Загрузка осуществляется средствами оркестратора Kubernetes по Helm-чарту:

```bash
helm repo add notifyhub {{ .helm_repo }}
helm upgrade --install notifyhub notifyhub/notifyhub \
  --version {{ .program_version }} \
  --namespace platform-notify \
  --create-namespace \
  --values values.prod.yaml
```

Контейнеры запускаются с привилегиями non-root, переменные окружения подгружаются из ConfigMap `notifyhub-config` и Secret `notifyhub-credentials`.

## 6. Входные данные

| Источник | Формат | Описание |
| --- | --- | --- |
| HTTP API `/v1/notify` | JSON | внешний публичный API публикации события |
| Очередь `notify.events` | JSON (AMQP) | внутренний канал публикации событий |
| Файл `subscriptions.csv` | CSV | пакетная загрузка подписок (адм. панель) |

Структура входного события (фрагмент JSON-схемы):

```json
{
  "type": "object",
  "required": ["event_id", "recipient_id", "channel", "payload"],
  "properties": {
    "event_id":      { "type": "string", "format": "uuid" },
    "recipient_id":  { "type": "string" },
    "channel":       { "enum": ["push", "sms", "email"] },
    "priority":      { "enum": ["low", "normal", "high"] },
    "payload":       { "type": "object" }
  }
}
```

## 7. Выходные данные

| Назначение | Формат | Куда выводится |
| --- | --- | --- |
| Команда отправки | JSON (AMQP) | очереди `routing.{push,sms,email}` |
| Журнал доставки | строки JSON | сток в Vector → Clickhouse |
| Метрики | OpenMetrics (Prometheus) | `:9090/metrics` |
| Отчёты | CSV/PDF | административная панель |

Среднесуточные объёмы:

- входящих событий — ≈ {{ .daily_events }} (среднее);
- сообщений в канал push — ≈ {{ .daily_push }};
- сообщений в канал SMS — ≈ {{ .daily_sms }};
- сообщений в канал email — ≈ {{ .daily_email }}.
