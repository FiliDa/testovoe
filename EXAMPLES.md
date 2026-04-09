# Демонстрация работы системы периодичности задач

## Примеры использования API

### 1. Ежедневная задача (каждые 3 дня)

```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "Ежедневный обзвон пациентов",
  "description": "Обзвонить пациентов для подтверждения визита",
  "status": "new",
  "periodicity_type": "daily",
  "periodicity_interval": 3
}
```

**Ответ:**
```json
{
  "id": 1,
  "title": "Ежедневный обзвон пациентов",
  "description": "Обзвонить пациентов для подтверждения визита",
  "status": "new",
  "periodicity_type": "daily",
  "periodicity_interval": 3,
  "next_occurrence": "2024-01-04T00:00:00Z"
}
```

### 2. Ежемесячная задача (1-е и 15-е число)

```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "Формирование отчетности",
  "description": "Подготовка ежемесячных отчетов",
  "status": "new",
  "periodicity_type": "monthly",
  "periodicity_days": [1, 15]
}
```

**Ответ:**
```json
{
  "id": 2,
  "title": "Формирование отчетности",
  "description": "Подготовка ежемесячных отчетов",
  "status": "new",
  "periodicity_type": "monthly",
  "periodicity_days": [1, 15],
  "next_occurrence": "2024-01-15T00:00:00Z"
}
```

### 3. Конкретные даты

```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "Плановый медосмотр",
  "description": "Ежегодный плановый медосмотр персонала",
  "status": "new",
  "periodicity_type": "specific_dates",
  "periodicity_dates": [
    "2024-03-15T00:00:00Z",
    "2024-06-15T00:00:00Z",
    "2024-09-15T00:00:00Z"
  ]
}
```

**Ответ:**
```json
{
  "id": 3,
  "title": "Плановый медосмотр",
  "description": "Ежегодный плановый медосмотр персонала",
  "status": "new",
  "periodicity_type": "specific_dates",
  "periodicity_dates": [
    "2024-03-15T00:00:00Z",
    "2024-06-15T00:00:00Z",
    "2024-09-15T00:00:00Z"
  ],
  "next_occurrence": "2024-03-15T00:00:00Z"
}
```

### 4. Четные дни месяца

```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "Инвентаризация на четные дни",
  "description": "Проверка медицинского оборудования",
  "status": "new",
  "periodicity_type": "even_odd",
  "periodicity_even_odd": "even"
}
```

**Ответ:**
```json
{
  "id": 4,
  "title": "Инвентаризация на четные дни",
  "description": "Проверка медицинского оборудования",
  "status": "new",
  "periodicity_type": "even_odd",
  "periodicity_even_odd": "even",
  "next_occurrence": "2024-01-02T00:00:00Z"
}
```

### 5. Обычная задача (без периодичности)

```json
POST /api/v1/tasks
Content-Type: application/json

{
  "title": "Разовая консультация",
  "description": "Консультация пациента по результатам анализов",
  "status": "new"
}
```

**Ответ:**
```json
{
  "id": 5,
  "title": "Разовая консультация",
  "description": "Консультация пациента по результатам анализов",
  "status": "new"
}
```

## Получение списка задач

```http
GET /api/v1/tasks
```

**Ответ:**
```json
[
  {
    "id": 1,
    "title": "Ежедневный обзвон пациентов",
    "description": "Обзвонить пациентов для подтверждения визита",
    "status": "new",
    "periodicity_type": "daily",
    "periodicity_interval": 3,
    "next_occurrence": "2024-01-04T00:00:00Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  },
  {
    "id": 2,
    "title": "Формирование отчетности",
    "description": "Подготовка ежемесячных отчетов",
    "status": "new",
    "periodicity_type": "monthly",
    "periodicity_days": [1, 15],
    "next_occurrence": "2024-01-15T00:00:00Z",
    "created_at": "2024-01-01T10:05:00Z",
    "updated_at": "2024-01-01T10:05:00Z"
  }
]
```

## Обновление задачи с периодичностью

```http
PUT /api/v1/tasks/1
Content-Type: application/json

{
  "title": "Ежедневный обзвон пациентов (обновлено)",
  "description": "Обзвон пациентов с новым скриптом",
  "status": "in_progress",
  "periodicity_type": "daily",
  "periodicity_interval": 2
}
```

**Ответ:**
```json
{
  "id": 1,
  "title": "Ежедневный обзвон пациентов (обновлено)",
  "description": "Обзвон пациентов с новым скриптом",
  "status": "in_progress",
  "periodicity_type": "daily",
  "periodicity_interval": 2,
  "next_occurrence": "2024-01-03T00:00:00Z",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T11:30:00Z"
}
```

## Как работает система

### Алгоритм вычисления следующих дат:

1. **Ежедневные задачи**: Текущая дата + interval дней
2. **Ежемесячные задачи**: Ближайший день в текущем месяце, либо первый день следующего месяца
3. **Конкретные даты**: Первая будущая дата из списка
4. **Четные/нечетные дни**: Ближайший подходящий день

### Валидация данных:
- Проверка корректности типа периодичности
- Валидация дней месяца (1-31)
- Проверка что конкретные даты находятся в будущем
- Валидация четности/нечетности

### Безопасность:
- Все запросы валидируются
- Защита от SQL-инъекций
- Rate limiting (100 запросов/мин)
- Security headers (CSP, XSS защита)