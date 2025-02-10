## Backup Tool
[![Go Report Card](https://goreportcard.com/badge/github.com/itocode21/backup-tool)](https://goreportcard.com/report/github.com/itocode21/Tender-backup-tool)
Its my solution for the   [backup-tool](https://roadmap.sh/projects/database-backup-utility) challenge from [roadmap.sh](https://roadmap.sh/).

## Описание

Backup Tool — это инструмент командной строки (CLI) для выполнения резервного копирования и восстановления данных из различных баз данных (MySQL, PostgreSQL, MongoDB). Инструмент поддерживает локальное хранение бэкапов и предоставляет удобный интерфейс для управления процессами.

## Требования
* Go 1.23+
* Docker (Для запуска через контейнеры)
* Установленные утилиты 
    ```mysqldump```, ```mysql``` (для MySQL)
    ```pg_dump```, ```psql``` (для PostgreSQL)
    ```mongodump```, ```mongorestore``` (для MongoDB)



## Установка

1.  Установите Go (версия 1.23 или выше). Вы можете скачать Go с [https://go.dev/dl/](https://go.dev/dl/).
2.  Установите Docker и Docker Compose. Docker Desktop можно скачать отсюда [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/).
3. Установите GolangCI-lint  с [https://github.com/golangci/golangci-lint#installation](https://github.com/golangci/golangci-lint#installation)
4.  Клонируйте репозиторий:

    ```bash
    git clone https://github.com/your-username/backup-tool.git
    cd backup-tool
    ``` 

    Установка Зависимостей:
    ```bash
    go mod tidy
    ```

    Соберите исполняемый файл:
    ```bash
    make build
    ```


## CLI-команды

1. Резервное копирование
```bash
   ./build/backup-tool --config pkg/config/mysql.yaml --type mysql --command backup --backup-file data/backups/mysql/mydb.sql
```

2. Восстановление
```bash
   ./build/backup-tool --config pkg/config/mysql.yaml --type mysql --command restore --backup-file data/backups/mysql/mydb.sql
```
3. Параметры CLI
```bash
--config: Путь к файлу конфигурации (обязательный).
--type: Тип базы данных (mysql, postgresql, mongodb) (обязательный).
--command: Команда для выполнения (backup, restore) (обязательный).
--backup-file: Путь к файлу бэкапа (для restore и backup).
```

## Пример файла конфигурации
```yaml
database:
  type: mysql
  host: localhost
  port: 3306
  username: root
  password: backup_password
  dbname: test_db

storage:
  local_path: data/backups
  cloud_type: s3
  bucket: my-backup-bucket

logging:
  level: info
  file: data/logs/app.log
  format: text

notification:
  slack_webhook_url: https://hooks.slack.com/services/...
  ```

# Со временем добавлю
1. Облачное хранилище
    * Поддержка загрузки бекапов в облачные хранилища(AWS S3, GCS, Yandex cloud)
2. Уведомления:
    * Реализую поддержку уведомлений через Slack.
3. Планировщик для автоматического выполнения бекапа
