#!/bin/bash

# Record start time
start_time=$(date +%s)

# Exit on error
set -e

echo "Setting up test database..."

# Check if test database exists and drop if it does
docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -c "SELECT 1 FROM pg_database WHERE datname = 'financial_test_db'" | grep -q 1 && \
  docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -c "DROP DATABASE financial_test_db"

# Create test database
docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -c "CREATE DATABASE financial_test_db"

echo "Running integration tests..."

# Run integration tests
go test -v

# Record end time and calculate duration
end_time=$(date +%s)
duration=$((end_time - start_time))
echo "Integration tests completed in ${duration} seconds!" 
