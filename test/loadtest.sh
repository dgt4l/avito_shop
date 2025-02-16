#!/bin/bash

if ! command -v vegeta &> /dev/null; then
    echo "Vegeta не установлен. Установите с помощью: go install github.com/tsenart/vegeta/v12@latest"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "Ошибка: Установите jq с помощью sudo apt-get install jq"
    exit 1
fi

# Конфигурация
BASE_URL="http://localhost:8080"
RATE="100/s"        # RPS
DURATION="30s"      # Длительность теста

USERB_CREDENTIALS='{"username": "abcde1234567", "password": "password1234"}'
USER_CREDENTIALS='{"username": "first1234567", "password": "password1234"}'
USER1_CREDENTIALS='{"username": "second1234567", "password": "password1234"}'

AUTH_ENDPOINT="/api/auth"
TEST_ENDPOINTS=(
    "/api/auth"
    "/api/info"
    "/api/sendCoin"
    "/api/?item=cup"
)

echo "Получаем 1 токен аутентификации..."
TOKENB=$(curl -s -X POST "$BASE_URL$AUTH_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "$USERB_CREDENTIALS" | jq -r '.token')

if [ -z "$TOKENB" ]; then
    echo "Ошибка: Не удалось получить токен. Проверьте данные аутентификации."
    exit 1
fi

echo "Получаем 2 токен аутентификации..."
TOKEN1=$(curl -s -X POST "$BASE_URL$AUTH_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "$USER_CREDENTIALS" | jq -r '.token')

if [ -z "$TOKEN1" ]; then
    echo "Ошибка: Не удалось получить токен. Проверьте данные аутентификации."
    exit 1
fi

echo "Получаем 3 токен аутентификации..."
TOKEN2=$(curl -s -X POST "$BASE_URL$AUTH_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "$USER1_CREDENTIALS" | jq -r '.token')

if [ -z "$TOKEN2" ]; then
    echo "Ошибка: Не удалось получить токен. Проверьте данные аутентификации."
    exit 1
fi

echo "Создаем тестового пользователя"
curl -s -X POST "$BASE_URL$AUTH_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "$USERB_CREDENTIALS"

echo "Создаем тестового пользователя"
curl -s -X POST "$BASE_URL$AUTH_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "$USER_CREDENTIALS"

echo "Создаем тестового пользователя"
curl -s -X POST "$BASE_URL$AUTH_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "$USER1_CREDENTIALS"

echo "Создаем файл целей..."
cat > targets.txt <<EOF
GET $BASE_URL/api/info
Authorization: Bearer $TOKEN1

GET $BASE_URL/api/info
Authorization: Bearer $TOKEN2

POST $BASE_URL/api/sendCoin
Authorization: Bearer $TOKEN2
Content-Type: application/json
@sendCoin_body2.json

POST $BASE_URL/api/sendCoin
Authorization: Bearer $TOKEN1
Content-Type: application/json
@sendCoin_body1.json

GET $BASE_URL/api/buy?item=cup
Authorization: Bearer $TOKENB
Content-Type: application/json

POST $BASE_URL/api/auth
Content-Type: application/json
@sendCoin_bodyb.json
EOF

# 3. Создаем тела запросов
echo '{"toUser": "second123456", "amount": 1}' > sendCoin_body1.json
echo '{"toUser": "first123456", "amount": 1}' > sendCoin_body2.json
echo '{"username": "abcde123456", "password": "password1234"}' > sendCoin_bodyb.json


# 4. Запуск теста
echo "Запускаем нагрузочный тест..."
vegeta attack \
  -rate="$RATE" \
  -duration="$DURATION" \
  -targets=targets.txt \
  > results.bin

# 5. Генерация отчетов
echo "Генерируем отчеты..."
vegeta report results.bin > report.txt
vegeta plot results.bin > plot.html
vegeta report -type=json results.bin > metrics.json

echo "Тестирование завершено!"
echo "Результаты:"
echo " - Текстовый отчет: report.txt"
echo " - Визуальный график: plot.html"
echo " - Метрики в JSON: metrics.json"

# Очистка временных файлов
rm -f targets.txt sendCoin_body2.json results.bin
rm -f targets.txt sendCoin_body1.json results.bin
rm -f targets.txt sendCoin_bodyb.json results.bin