#!/bin/bash

# Navigate to the project directory
cd /home/market-league || exit

# Pull the latest changes from the repository
git pull origin main

# Run Docker Compose to rebuild and restart the services
docker compose -f docker-compose.prod.yml up --build -d