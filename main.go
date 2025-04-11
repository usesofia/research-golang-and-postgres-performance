package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=financial_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Tag{}, &FinancialRecord{})

	// Apply database indexes
	ApplyIndexes(db)

	// Initialize router
	r := gin.Default()

	// Routes
	r.POST("/organizations/:organizationId/tags", createTag(db))
	r.GET("/organizations/:organizationId/tags", listTags(db))
	r.POST("/organizations/:organizationId/financial-records", createFinancialRecord(db))
	r.POST("/organizations/:organizationId/financial-records/bulk", createFinancialRecordsBulk(db))
	r.GET("/organizations/:organizationId/financial-records", listFinancialRecords(db))
	r.GET("/organizations/:organizationId/financial-records/reports/cash-flow", getCashFlowReport(db))

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
