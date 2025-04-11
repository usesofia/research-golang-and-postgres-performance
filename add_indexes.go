package main

import (
	"log"

	"gorm.io/gorm"
)

// ApplyIndexes creates database indexes to optimize queries
func ApplyIndexes(db *gorm.DB) {
	log.Println("Applying database indexes...")

	// Index for cash flow report query
	err := db.Exec("CREATE INDEX IF NOT EXISTS idx_financial_records_org_date ON financial_records (organization_id, due_date)").Error
	if err != nil {
		log.Printf("Warning: Failed to create org_date index: %v", err)
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_financial_records_due_date ON financial_records (due_date)").Error
	if err != nil {
		log.Printf("Warning: Failed to create due_date index: %v", err)
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_financial_records_direction ON financial_records (direction)").Error
	if err != nil {
		log.Printf("Warning: Failed to create direction index: %v", err)
	}

	// Indexes for financial_record_tags join table
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_financial_record_tags_record_id ON financial_record_tags (financial_record_id)").Error
	if err != nil {
		log.Printf("Warning: Failed to create financial_record_tags_record_id index: %v", err)
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_financial_record_tags_tag_id ON financial_record_tags (tag_id)").Error
	if err != nil {
		log.Printf("Warning: Failed to create financial_record_tags_tag_id index: %v", err)
	}

	// Apply indexes for all the new models
	addModelIndexes(db)

	log.Println("Database indexes applied successfully")
}

// addModelIndexes adds indexes for all new domain models
func addModelIndexes(db *gorm.DB) {
	// E-commerce Domain
	createIndex(db, "products", "category_id")
	createIndex(db, "products", "created_at")
	createIndex(db, "orders", "customer_id")
	createIndex(db, "orders", "status")
	createIndex(db, "orders", "order_date")
	createIndex(db, "order_items", "order_id")
	createIndex(db, "order_items", "product_id")
	
	// Healthcare Domain
	createIndex(db, "patients", "created_at")
	createIndex(db, "appointments", "patient_id")
	createIndex(db, "appointments", "doctor_id")
	createIndex(db, "appointments", "start_time")
	createIndex(db, "appointments", "status")
	createIndex(db, "prescriptions", "patient_id")
	
	// Education Domain
	createIndex(db, "students", "created_at")
	createIndex(db, "courses", "teacher_id")
	createIndex(db, "courses", "department")
	createIndex(db, "grades", "student_id")
	createIndex(db, "grades", "course_id")
	createIndex(db, "assignments", "course_id")
	createIndex(db, "assignments", "due_date")
	
	// Social Media Domain
	createIndex(db, "users", "username")
	createIndex(db, "users", "created_at")
	createIndex(db, "posts", "user_id")
	createIndex(db, "posts", "post_date")
	createIndex(db, "comments", "post_id")
	createIndex(db, "comments", "user_id")
	createIndex(db, "follows", "follower_id")
	createIndex(db, "follows", "followee_id")
	
	// Human Resources Domain
	createIndex(db, "employees", "department_id")
	createIndex(db, "employees", "position_id")
	createIndex(db, "employees", "hire_date")
	createIndex(db, "leaves", "employee_id")
	createIndex(db, "leaves", "start_date")
	createIndex(db, "leaves", "status")
	createIndex(db, "evaluations", "employee_id")
	
	// Content Management Domain
	createIndex(db, "websites", "owner_id")
	createIndex(db, "pages", "website_id")
	createIndex(db, "pages", "slug")
	createIndex(db, "pages", "published")
	
	// IoT Domain
	createIndex(db, "devices", "home_id")
	createIndex(db, "devices", "status")
	createIndex(db, "device_readings", "device_id")
	createIndex(db, "device_readings", "reading_type")
	createIndex(db, "device_readings", "timestamp")
	
	// Fitness & Wellness Domain
	createIndex(db, "workouts", "user_id")
	createIndex(db, "workouts", "start_time")
	createIndex(db, "workout_exercises", "workout_id")
	createIndex(db, "meal_plans", "user_id")
	createIndex(db, "meals", "meal_plan_id")
	
	// Finance Domain
	createIndex(db, "budgets", "organization_id")
	createIndex(db, "budget_categories", "budget_id")
	createIndex(db, "invoices", "organization_id")
	createIndex(db, "invoices", "due_date")
	createIndex(db, "invoices", "status")
	
	// Real Estate Domain
	createIndex(db, "properties", "owner_id")
	createIndex(db, "properties", "is_for_sale")
	createIndex(db, "properties", "is_for_rent")
	createIndex(db, "listings", "property_id")
	createIndex(db, "listings", "agent_id")
	createIndex(db, "listings", "status")
	
	// Hospitality Domain
	createIndex(db, "hotels", "rating")
	createIndex(db, "hotel_rooms", "hotel_id")
	createIndex(db, "hotel_rooms", "room_type")
	createIndex(db, "bookings", "hotel_id")
	createIndex(db, "bookings", "room_id")
	createIndex(db, "bookings", "guest_id")
	createIndex(db, "bookings", "check_in_date")
	createIndex(db, "bookings", "status")
	
	// Event Management Domain
	createIndex(db, "events", "organizer_id")
	createIndex(db, "events", "start_date")
	createIndex(db, "events", "status")
	createIndex(db, "tickets", "event_id")
	createIndex(db, "tickets", "attendee_id")
	createIndex(db, "event_sessions", "event_id")
	
	// Logistics Domain
	createIndex(db, "shipments", "carrier_id")
	createIndex(db, "shipments", "status")
	createIndex(db, "shipments", "ship_date")
	createIndex(db, "packages", "shipment_id")
	
	// Project Management Domain
	createIndex(db, "projects", "manager_id")
	createIndex(db, "projects", "status")
	createIndex(db, "tasks", "project_id")
	createIndex(db, "tasks", "assignee_id")
	createIndex(db, "tasks", "status")
	createIndex(db, "tasks", "due_date")
	createIndex(db, "milestones", "project_id")
	
	// Legal Domain
	createIndex(db, "contracts", "party_one_id")
	createIndex(db, "contracts", "party_two_id")
	createIndex(db, "contracts", "status")
	createIndex(db, "clauses", "contract_id")
	
	// Agriculture Domain
	createIndex(db, "farms", "owner_id")
	createIndex(db, "fields", "farm_id")
	createIndex(db, "livestock", "farm_id")
}

// createIndex creates a single index with error handling
func createIndex(db *gorm.DB, table string, column string) {
	indexName := "idx_" + table + "_" + column
	err := db.Exec("CREATE INDEX IF NOT EXISTS " + indexName + " ON " + table + " (" + column + ")").Error
	if err != nil {
		log.Printf("Warning: Failed to create %s index: %v", indexName, err)
	}
} 