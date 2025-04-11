#!/bin/bash

# Default values
TEST_NUMBER=${1:-0}
DURATION=${2:-60s}

echo "Running test ${TEST_NUMBER} for ${DURATION}..."

# Clearing the db
echo "Deleting financial_record_tags..."
docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -d financial_db -c "DELETE FROM financial_record_tags;"
echo "Deleting tags..."
docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -d financial_db -c "DELETE FROM tags;"
echo "Deleting financial_records..."
docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -d financial_db -c "DELETE FROM financial_records;"

echo "Running populate.js..."
K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_EXPORT=./reports/test-${TEST_NUMBER}-populate.html k6 run --vus 100 --duration ${DURATION} populate.js

# Connect to the database and getting count of tags
docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -d financial_db -c "SELECT COUNT(*) FROM tags;" > ./reports/test-${TEST_NUMBER}-populate-tags-count.txt
echo "Tags:\n $(cat ./reports/test-${TEST_NUMBER}-populate-tags-count.txt)"

# Connect to the database and getting count of financial records
docker exec -it research-golang-and-postgres-performance-db-1 psql -U postgres -d financial_db -c "SELECT COUNT(*) FROM financial_records;" > ./reports/test-${TEST_NUMBER}-populate-financial-records-count.txt
echo "Financial records:\n $(cat ./reports/test-${TEST_NUMBER}-populate-financial-records-count.txt)"
