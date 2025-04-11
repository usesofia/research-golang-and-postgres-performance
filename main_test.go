package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	// Setup test environment
	gin.SetMode(gin.TestMode)
	
	// Use a test database
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=financial_test_db port=5432 sslmode=disable"
	}

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	// Migrate the schema
	testDB.AutoMigrate(&Tag{}, &FinancialRecord{})
	
	// Apply database indexes
	ApplyIndexes(testDB)
	
	// Setup router with routes
	router = gin.Default()
	router.POST("/organizations/:organizationId/tags", createTag(testDB))
	router.GET("/organizations/:organizationId/tags", listTags(testDB))
	router.POST("/organizations/:organizationId/financial-records", createFinancialRecord(testDB))
	router.POST("/organizations/:organizationId/financial-records/bulk", createFinancialRecordsBulk(testDB))
	router.GET("/organizations/:organizationId/financial-records", listFinancialRecords(testDB))
	router.GET("/organizations/:organizationId/financial-records/reports/cash-flow", getCashFlowReport(testDB))
	
	// Run tests
	exitCode := m.Run()
	
	// Clean up database after tests
	cleanupTestDB()
	
	os.Exit(exitCode)
}

func cleanupTestDB() {
	// Drop all tables
	testDB.Exec("DROP TABLE IF EXISTS financial_record_tags CASCADE")
	testDB.Exec("DROP TABLE IF EXISTS financial_records CASCADE")
	testDB.Exec("DROP TABLE IF EXISTS tags CASCADE")
}

func clearTables() {
	testDB.Exec("DELETE FROM financial_record_tags")
	testDB.Exec("DELETE FROM financial_records")
	testDB.Exec("DELETE FROM tags")
}

func TestCreateTag(t *testing.T) {
	clearTables()
	
	// Create test data
	tag := map[string]interface{}{
		"name": "Test Tag",
	}
	jsonData, _ := json.Marshal(tag)
	
	// Create request
	req := httptest.NewRequest("POST", "/organizations/1/tags", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve the request
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Parse response
	var response Tag
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	
	// Validate response
	assert.Equal(t, "Test Tag", response.Name)
	assert.Equal(t, uint(1), response.OrganizationID)
	assert.NotZero(t, response.ID)
}

func TestListTags(t *testing.T) {
	clearTables()
	
	// Create test data
	testDB.Create(&Tag{Name: "Tag 1", OrganizationID: 1})
	testDB.Create(&Tag{Name: "Tag 2", OrganizationID: 1})
	
	// Create request
	req := httptest.NewRequest("GET", "/organizations/1/tags", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve the request
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response struct {
		Data []Tag `json:"data"`
		Pagination struct {
			CurrentPage int   `json:"current_page"`
			PageSize    int   `json:"page_size"`
			TotalItems  int64 `json:"total_items"`
			TotalPages  int64 `json:"total_pages"`
		} `json:"pagination"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	
	// Validate response
	assert.Len(t, response.Data, 2)
	assert.Equal(t, "Tag 1", response.Data[0].Name)
	assert.Equal(t, "Tag 2", response.Data[1].Name)
	
	// Validate pagination
	assert.Equal(t, 1, response.Pagination.CurrentPage)
	assert.Equal(t, 20, response.Pagination.PageSize)
	assert.Equal(t, int64(2), response.Pagination.TotalItems)
	assert.Equal(t, int64(1), response.Pagination.TotalPages)
}

func TestCreateFinancialRecord(t *testing.T) {
	clearTables()
	
	// Create test data - first create a tag
	tag := Tag{Name: "Expense Tag", OrganizationID: 1}
	testDB.Create(&tag)
	
	// Create financial record with tag
	record := map[string]interface{}{
		"direction": "OUT",
		"amount":    100.50,
		"dueDate":   time.Now().Format(time.RFC3339),
		"tags":      []map[string]interface{}{{"id": tag.ID}},
	}
	jsonData, _ := json.Marshal(record)
	
	// Create request
	req := httptest.NewRequest("POST", "/organizations/1/financial-records", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve the request
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Parse response
	var response FinancialRecord
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	
	// Validate response
	assert.Equal(t, "OUT", response.Direction)
	assert.Equal(t, 100.50, response.Amount)
	assert.Equal(t, uint(1), response.OrganizationID)
}

func TestCreateFinancialRecordsBulk(t *testing.T) {
	clearTables()
	
	// Create test data - first create tags
	tag1 := Tag{Name: "Income Tag", OrganizationID: 1}
	tag2 := Tag{Name: "Expense Tag", OrganizationID: 1}
	testDB.Create(&tag1)
	testDB.Create(&tag2)
	
	// Create financial records with tags
	records := []map[string]interface{}{
		{
			"direction": "IN",
			"amount":    1000.0,
			"dueDate":   time.Now().Format(time.RFC3339),
			"tags":      []map[string]interface{}{{"id": tag1.ID}},
		},
		{
			"direction": "OUT",
			"amount":    500.0,
			"dueDate":   time.Now().Format(time.RFC3339),
			"tags":      []map[string]interface{}{{"id": tag2.ID}},
		},
	}
	jsonData, _ := json.Marshal(records)
	
	// Create request
	req := httptest.NewRequest("POST", "/organizations/1/financial-records/bulk", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve the request
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Parse response
	var response []FinancialRecord
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	
	// Validate response
	assert.Len(t, response, 2)
	assert.Equal(t, "IN", response[0].Direction)
	assert.Equal(t, "OUT", response[1].Direction)
}

func TestListFinancialRecords(t *testing.T) {
	clearTables()
	
	// Create test data
	tag := Tag{Name: "Test Tag", OrganizationID: 1}
	testDB.Create(&tag)
	
	record1 := FinancialRecord{
		Direction:      "IN",
		Amount:         1000.0,
		DueDate:        time.Now(),
		OrganizationID: 1,
	}
	testDB.Create(&record1)
	testDB.Exec("INSERT INTO financial_record_tags (financial_record_id, tag_id) VALUES (?, ?)", record1.ID, tag.ID)
	
	record2 := FinancialRecord{
		Direction:      "OUT",
		Amount:         500.0,
		DueDate:        time.Now(),
		OrganizationID: 1,
	}
	testDB.Create(&record2)
	
	// Create request
	req := httptest.NewRequest("GET", "/organizations/1/financial-records?page=1&page_size=10", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve the request
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	
	// Validate response
	records, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, records, 2)
	
	pagination, ok := response["pagination"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(1), pagination["current_page"])
	assert.Equal(t, float64(2), pagination["total_items"])
}

func TestCashFlowReport(t *testing.T) {
	clearTables()
	
	// Create test data
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	
	// Create records for current month
	testDB.Create(&FinancialRecord{
		Direction:      "IN",
		Amount:         2000.0,
		DueDate:        now,
		OrganizationID: 1,
	})
	testDB.Create(&FinancialRecord{
		Direction:      "OUT",
		Amount:         1000.0,
		DueDate:        now,
		OrganizationID: 1,
	})
	
	// Create records for last month
	testDB.Create(&FinancialRecord{
		Direction:      "IN",
		Amount:         1500.0,
		DueDate:        lastMonth,
		OrganizationID: 1,
	})
	testDB.Create(&FinancialRecord{
		Direction:      "OUT",
		Amount:         800.0,
		DueDate:        lastMonth,
		OrganizationID: 1,
	})
	
	// Create request
	req := httptest.NewRequest("GET", "/organizations/1/financial-records/reports/cash-flow", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve the request
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response CashFlowReport
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	
	// Validate response has monthly data
	assert.GreaterOrEqual(t, len(response.MonthlyData), 2)
	
	// Find current month and last month data
	currentYearMonth := fmt.Sprintf("%d-%d", now.Year(), int(now.Month()))
	lastYearMonth := fmt.Sprintf("%d-%d", lastMonth.Year(), int(lastMonth.Month()))
	
	foundCurrent := false
	foundLast := false

	for _, data := range response.MonthlyData {
		yearMonth := fmt.Sprintf("%d-%d", data.Year, data.Month)
		if yearMonth == currentYearMonth {
			foundCurrent = true
			assert.InDelta(t, 2000.0, data.In, 0.01)
			assert.InDelta(t, 1000.0, data.Out, 0.01)
		} else if yearMonth == lastYearMonth {
			foundLast = true
			assert.InDelta(t, 1500.0, data.In, 0.01)
			assert.InDelta(t, 800.0, data.Out, 0.01)
		}
	}
	
	assert.True(t, foundCurrent || foundLast, "Should find data for current month or last month")
} 