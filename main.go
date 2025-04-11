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
	db.AutoMigrate(&Tag{}, &FinancialRecord{},
		// E-commerce Domain
		&Product{}, &Category{}, &Customer{}, &Order{}, &OrderItem{}, &Image{}, &Address{}, &Transaction{},
		// Healthcare Domain
		&Patient{}, &Doctor{}, &Appointment{}, &Prescription{}, &Medication{},
		// Education Domain
		&Student{}, &Course{}, &Teacher{}, &Grade{}, &Assignment{},
		// Transportation Domain
		&Vehicle{}, &Driver{}, &Trip{}, &Passenger{},
		// Social Media Domain
		&User{}, &Post{}, &Comment{}, &Follow{}, &Attachment{},
		// Human Resources Domain
		&Employee{}, &Department{}, &Position{}, &Leave{}, &Evaluation{},
		// Content Management Domain
		&Website{}, &Page{}, &MenuItem{},
		// IoT & Smart Home Domain
		&Device{}, &Home{}, &Room{}, &DeviceReading{},
		// Fitness & Wellness Domain
		&Workout{}, &Exercise{}, &WorkoutExercise{}, &MealPlan{}, &Meal{}, &MealFood{},
		// Finance Domain
		&Budget{}, &BudgetCategory{}, &Invoice{}, &InvoiceItem{},
		// Real Estate Domain
		&Property{}, &Listing{}, &Agent{},
		// Hospitality Domain
		&Hotel{}, &HotelRoom{}, &Booking{}, &Guest{},
		// Event Management Domain
		&Event{}, &Ticket{}, &EventSession{}, &Speaker{},
		// Manufacturing Domain
		&Component{}, &Supplier{},
		// Logistics Domain
		&Shipment{}, &Package{}, &Carrier{},
		// Project Management Domain
		&Project{}, &Task{}, &Milestone{}, &TeamMember{},
		// Legal Domain
		&Contract{}, &Clause{}, &LegalEntity{},
		// Agriculture Domain
		&Farm{}, &Field{}, &Livestock{})
	
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