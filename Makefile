# Переменные
GO := go
BINARY_NAME := backup-tool
BUILD_DIR := build
CONFIG_DIR := pkg/config
DATA_DIR := data
LOGS_DIR := $(DATA_DIR)/logs
BACKUPS_DIR := $(DATA_DIR)/backups

# Цели (targets)
.PHONY: all clean build run test docker-build docker-run

# Сборка приложения
all: clean build

# Очистка временных файлов и директорий
clean:
    @echo "Cleaning up..."
    rm -rf $(BUILD_DIR)
    rm -rf $(LOGS_DIR)/*
    rm -rf $(BACKUPS_DIR)/*

# Создание директорий для сборки
prepare:
    @mkdir -p $(BUILD_DIR)
    @mkdir -p $(LOGS_DIR)
    @mkdir -p $(BACKUPS_DIR)

# Сборка исполняемого файла
build: prepare
    @echo "Building the application..."
    $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
    @echo "Build completed. Binary is located in $(BUILD_DIR)/$(BINARY_NAME)"

# Запуск приложения
run: build
    @echo "Running the application..."
    $(BUILD_DIR)/$(BINARY_NAME) \
        --config $(CONFIG_DIR)/test_config_mysql.yaml \
        --type mysql \
        --command backup \
        --backup-file $(BACKUPS_DIR)/mysql/mydb.sql

# Запуск тестов
test:
    @echo "Running tests..."
    $(GO) test ./... -v

# Сборка Docker-образа
docker-build:
    @echo "Building Docker image..."
    docker-compose build

# Запуск Docker-контейнеров
docker-run:
    @echo "Running Docker containers..."
    docker-compose up

# Остановка Docker-контейнеров
docker-stop:
    @echo "Stopping Docker containers..."
    docker-compose down

# Удаление Docker-контейнеров и томов
docker-clean:
    @echo "Cleaning up Docker resources..."
    docker-compose down --volumes
    docker system prune -f