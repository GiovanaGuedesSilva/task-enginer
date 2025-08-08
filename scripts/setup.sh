#!/bin/bash

# Setup script for Task Engine project

echo "Creating bin directory..."
mkdir -p ./bin

if [ ! -f .env ]; then
  echo "Copying .env.example to .env..."
  if [ -f .env.example ]; then
    cp .env.example .env
  else
    echo ".env.example not found, skipping"
  fi
else
  echo ".env already exists, skipping"
fi

echo "Downloading Go module dependencies..."
go mod download

echo "Tidying Go modules..."
go mod tidy

echo "Running database migrations..."
./scripts/run-migrations.sh up

echo "Setup complete!"
