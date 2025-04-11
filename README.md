# Financial Records API

A REST API built with Go and PostgreSQL for managing financial records and tags.

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher

## Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE financial_db;
```

2. Set up the environment variables (optional):
```bash
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=financial_db port=5432 sslmode=disable"
```

3. Install dependencies:
```bash
go mod tidy
```

4. Run the application:
```bash
go run .
```

The server will start on port 8080.

## API Endpoints

### Create a Tag
```
POST /organizations/:organizationId/tags
```
Request body:
```json
{
    "name": "string"
}
```

### Create a Financial Record
```
POST /organizations/:organizationId/financial-records
```
Request body:
```json
{
    "direction": "IN|OUT",
    "amount": number,
    "tags": [tag_ids],
    "dueDate": "YYYY-MM-DD"
}
```

### List Financial Records
```
GET /organizations/:organizationId/financial-records?tags=1,2,3
```

### Get Cash Flow Report
```
GET /organizations/:organizationId/financial-records/reports/cash-flow
```
Returns monthly cash flow data for the last two years.

## Running Tests

### Integration Tests

To run the integration tests, execute the provided script:

```bash
./scripts/run-integration-tests.sh
```

This script will:
1. Create a test database (`financial_test_db`)
2. Run all integration tests
3. Show detailed test results

The integration tests verify all API endpoints, including:
- Creating and listing tags
- Creating individual and bulk financial records 
- Listing financial records with pagination
- Generating cash flow reports

### Test Requirements

- PostgreSQL should be running locally with default settings
- The postgres user should have permission to create databases
- Go testing dependencies will be automatically installed
