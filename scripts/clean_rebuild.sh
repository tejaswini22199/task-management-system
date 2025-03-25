#!/bin/bash

echo "🚀 Stopping and removing all containers..."
CONTAINER_IDS=$(docker ps -aq)

if [ -n "$CONTAINER_IDS" ]; then
    docker stop $CONTAINER_IDS
    docker rm $CONTAINER_IDS
    echo "✅ All containers stopped and removed."
else
    echo "⚠️ No containers found."
fi

echo "🗑️ Removing all images containing 'task' in the repository name..."
IMAGE_IDS=$(docker images | grep "task" | awk '{print $3}')

if [ -n "$IMAGE_IDS" ]; then
    docker rmi -f $IMAGE_IDS
    echo "✅ All task-related images removed."
else
    echo "⚠️ No matching images found."
fi

echo "🛠️ Removing Docker volumes..."
docker volume prune -f

echo "🔥 Building the project from scratch..."
docker build --no-cache -t task-management-system .

echo "📦 Bringing down all Docker services..."
docker-compose down -v

echo "🚀 Starting services in detached mode..."
docker-compose up -d

echo "🎉 Rebuild complete! Everything is fresh & running!"
