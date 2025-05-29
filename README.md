# сервис для хранения и управления цитатами

сервис может быть запущен как локально, так и через docker.

## сборка docker-образа
```bash
docker build -t quotes-service .
```

## запуск через docker
```bash
docker-compose up -d
```

## запуск локально
```bash
# Сборка
go build -o quotes-service ./cmd/server

# Запуск
./quotes-service
```

для запуска на другом порту:
```bash
PORT=8081 ./quotes-service
```

## примеры использования

### добавление новой цитаты
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"author":"Линус Торвальдс", "quote":"Говорить, что ваш код самодокументируемый — это как говорить, что ваш дом самоубираемый."}' \
  http://localhost:8082/quotes
```

### получение всех цитат
```bash
curl http://localhost:8082/quotes
```

### получение случайной цитаты
```bash
curl http://localhost:8082/quotes/random
```

### фильтрация по автору
```bash
curl 'http://localhost:8082/quotes?author=Линус Торвальдс'
```

### удаление цитаты по ID
```bash
curl -X DELETE http://localhost:8082/quotes/1
```
