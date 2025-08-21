#!/bin/bash

# ErgrAcer API - Container Setup Script
# This script creates PostgreSQL or API containers with the same settings as docker-compose.yml
# Usage: ./setup-postgres.sh [db|api]

set -e

# Check for required argument
if [ $# -eq 0 ]; then
    echo "Error: Missing required argument"
    echo "Usage: $0 [db|api]"
    echo "  db  - Setup PostgreSQL database container"
    echo "  api - Setup ErgrAcer API container"
    exit 1
fi

SERVICE_TYPE="$1"

# Validate argument
if [ "$SERVICE_TYPE" != "db" ] && [ "$SERVICE_TYPE" != "api" ]; then
    echo "Error: Invalid argument '$SERVICE_TYPE'"
    echo "Usage: $0 [db|api]"
    echo "  db  - Setup PostgreSQL database container"
    echo "  api - Setup ErgrAcer API container"
    exit 1
fi

# Configuration variables
if [ "$SERVICE_TYPE" = "db" ]; then
    CONTAINER_NAME="ergracer-postgres"
    POSTGRES_VERSION="17"
    POSTGRES_DB="ergracer"
    POSTGRES_USER="ergracer"
    POSTGRES_PASSWORD="password123"
    POSTGRES_PORT="5432"
    VOLUME_NAME="ergracer-postgres-data"
else
    CONTAINER_NAME="ergracer-api"
    API_PORT="8080"
    API_IMAGE="ergracer-api"
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Stop and remove existing container if it exists
cleanup_existing() {
    if docker ps -a --format 'table {{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        print_warning "Existing container '${CONTAINER_NAME}' found."
        read -p "Do you want to remove it? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Stopping and removing existing container..."
            docker stop "${CONTAINER_NAME}" 2>/dev/null || true
            docker rm "${CONTAINER_NAME}" 2>/dev/null || true
            print_success "Existing container removed."
        else
            print_error "Cannot proceed with existing container. Exiting."
            exit 1
        fi
    fi
}

# Create Docker volume if it doesn't exist (only for database)
create_volume() {
    if [ "$SERVICE_TYPE" = "db" ]; then
        if ! docker volume ls | grep -q "${VOLUME_NAME}"; then
            print_status "Creating Docker volume '${VOLUME_NAME}'..."
            docker volume create "${VOLUME_NAME}"
            print_success "Volume created."
        else
            print_status "Volume '${VOLUME_NAME}' already exists."
        fi
    fi
}

# Start container based on service type
start_container() {
    if [ "$SERVICE_TYPE" = "db" ]; then
        start_postgres_container
    else
        start_api_container
    fi
}

# Start PostgreSQL container
start_postgres_container() {
    print_status "Starting PostgreSQL ${POSTGRES_VERSION} container..."
    
    docker run -d \
        --name "${CONTAINER_NAME}" \
        --restart unless-stopped \
        -e POSTGRES_DB="${POSTGRES_DB}" \
        -e POSTGRES_USER="${POSTGRES_USER}" \
        -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
        -e POSTGRES_INITDB_ARGS="--auth-host=scram-sha-256" \
        -p "${POSTGRES_PORT}:5432" \
        -v "${VOLUME_NAME}:/var/lib/postgresql/data" \
        --health-cmd="pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}" \
        --health-interval=10s \
        --health-timeout=5s \
        --health-retries=5 \
        "postgres:${POSTGRES_VERSION}"
    
    print_success "PostgreSQL container started."
}

# Start API container
start_api_container() {
    print_status "Starting ErgrAcer API container..."
    
    # Check if config.yaml exists
    if [ ! -f "config.yaml" ]; then
        print_error "config.yaml not found. Please copy config.yaml.example to config.yaml and configure it."
        exit 1
    fi
    
    # Check if API image exists
    if ! docker images | grep -q "${API_IMAGE}"; then
        print_warning "API image '${API_IMAGE}' not found. Building it now..."
        docker build -t "${API_IMAGE}" .
    fi
    
    docker run -d \
        --name "${CONTAINER_NAME}" \
        --restart unless-stopped \
        -p "${API_PORT}:8080" \
        -v "$(pwd)/config.yaml:/app/config.yaml:ro" \
        --health-cmd="wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1" \
        --health-interval=30s \
        --health-timeout=10s \
        --health-retries=3 \
        "${API_IMAGE}"
    
    print_success "ErgrAcer API container started."
}

# Wait for container to be ready
wait_for_container() {
    if [ "$SERVICE_TYPE" = "db" ]; then
        wait_for_postgres
    else
        wait_for_api
    fi
}

# Wait for PostgreSQL to be ready
wait_for_postgres() {
    print_status "Waiting for PostgreSQL to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker exec "${CONTAINER_NAME}" pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" > /dev/null 2>&1; then
            print_success "PostgreSQL is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 2
        ((attempt++))
    done
    
    echo
    print_error "PostgreSQL failed to start within expected time."
    print_status "Check container logs with: docker logs ${CONTAINER_NAME}"
    return 1
}

# Wait for API to be ready
wait_for_api() {
    print_status "Waiting for ErgrAcer API to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -o /dev/null -w "%{http_code}" "http://localhost:${API_PORT}/health" | grep -q "200"; then
            print_success "ErgrAcer API is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 2
        ((attempt++))
    done
    
    echo
    print_error "ErgrAcer API failed to start within expected time."
    print_status "Check container logs with: docker logs ${CONTAINER_NAME}"
    return 1
}

# Display connection information
show_connection_info() {
    echo
    if [ "$SERVICE_TYPE" = "db" ]; then
        show_postgres_info
    else
        show_api_info
    fi
}

# Display PostgreSQL connection information
show_postgres_info() {
    print_success "PostgreSQL container is running!"
    echo
    echo "Connection Details:"
    echo "  Host: localhost"
    echo "  Port: ${POSTGRES_PORT}"
    echo "  Database: ${POSTGRES_DB}"
    echo "  Username: ${POSTGRES_USER}"
    echo "  Password: ${POSTGRES_PASSWORD}"
    echo
    echo "Connection URL:"
    echo "  DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
    echo
    echo "Useful Commands:"
    echo "  Connect to database: docker exec -it ${CONTAINER_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}"
    echo "  View logs: docker logs ${CONTAINER_NAME}"
    echo "  Stop container: docker stop ${CONTAINER_NAME}"
    echo "  Start container: docker start ${CONTAINER_NAME}"
    echo "  Remove container: docker rm ${CONTAINER_NAME}"
    echo "  Remove volume: docker volume rm ${VOLUME_NAME}"
}

# Display API connection information
show_api_info() {
    print_success "ErgrAcer API container is running!"
    echo
    echo "API Details:"
    echo "  URL: http://localhost:${API_PORT}"
    echo "  Health Check: http://localhost:${API_PORT}/health"
    echo "  API Docs: http://localhost:${API_PORT}/api/v1"
    echo
    echo "Useful Commands:"
    echo "  View logs: docker logs ${CONTAINER_NAME}"
    echo "  Follow logs: docker logs -f ${CONTAINER_NAME}"
    echo "  Stop container: docker stop ${CONTAINER_NAME}"
    echo "  Start container: docker start ${CONTAINER_NAME}"
    echo "  Remove container: docker rm ${CONTAINER_NAME}"
    echo "  Rebuild image: docker build -t ${API_IMAGE} ."
}

# Main execution
main() {
    if [ "$SERVICE_TYPE" = "db" ]; then
        echo "ErgrAcer API - PostgreSQL Setup"
    else
        echo "ErgrAcer API - Container Setup"
    fi
    echo "==============================="
    echo
    
    check_docker
    cleanup_existing
    create_volume
    start_container
    
    if wait_for_container; then
        show_connection_info
    else
        exit 1
    fi
}

# Run main function
main "$@"