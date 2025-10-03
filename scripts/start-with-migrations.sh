#!/bin/sh

echo "Waiting for database to be ready..."
until pg_isready -h db -p 5432 -U postgres; do
  echo "Database is unavailable - sleeping"
  sleep 1
done

echo "Database is ready - running migrations..."
./migrationsfa -mig up 

echo "Migrations completed - starting main application..."
./mainaa

