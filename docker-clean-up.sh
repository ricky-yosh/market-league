#!/bin/bash
# Remove All Containers:
docker rm -f $(docker ps -a -q)

# Remove All Images:
docker rmi -f $(docker images -a -q)

# Remove All Networks:
docker network rm $(docker network ls -q)

# Remove All Volumes:
# docker volume rm $(docker volume ls -q)

# Remove Unused Docker Resources:
# docker system prune -a