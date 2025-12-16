#!/bin/bash
set -e

# Ждём, пока БД станет доступна
until PGPASSWORD="$POSTGRES_PASSWORD" psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  echo "⏳ Ожидание PostgreSQL..."
  sleep 2
done

echo "PostgreSQL готов. Создаём тестовые данные (если ещё не созданы)..."

# Выполняем идемпотентный SQL-скрипт
PGPASSWORD="$POSTGRES_PASSWORD" psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /testdata/testdata.sql

echo "Тестовые данные готовы."