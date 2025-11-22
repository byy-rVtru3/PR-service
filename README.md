```bash
# Запустить сервис
docker-compose up --build

# Сервис будет доступен на http://localhost:8080
```

## Makefile команды

```bash
make build    # Собрать Docker образы
make up       # Запустить сервис
make down     # Остановить сервис
make restart  # Перезапустить сервис
make logs     # Показать логи
make clean    # Остановить и удалить вольюмы
```