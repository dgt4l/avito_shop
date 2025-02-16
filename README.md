## Конфигурация проекта

Для демонстрации работы используются тестовые конфигурационные файлы:
`deploy/.env.dev`
`test_config.yaml`

Конфигурация линтера:
`.golangci.yml`

## Запуск проекта

### Варианты запуска

#### Makefile

С помощью команды 

```shell
make docker-buildup 
```

#### Docker

```shell
cd deploy
```
Включить сборку:
```shell
docker compose --env-file .env.dev up --build -d
```
Без сборки:
```shell
docker compose --env-file .env.dev up -d
```
#### Локально(без поднятия БД)

```shell
go run cmd/avito_shop/main.go
```

## Makefile

Помимо запуска Makefile различный функционал. К примеру запуск тестов:

```shell
make test
```

Также можно получить процент покрытия тестами (процент покрытия, включая директории `dto` и `models` составляет 41,6%):

```shell
make coverage
```

Полный функционал с описанием можно увидеть после ввода команды:

```shell
make help
```

## Нагрузочное тестирование

Результаты нагрузочного тестирования находятся в `test/`, файлы report.txt, metrics.json, plot.html. Сам скрипт для запуска - `loadtest.sh`


