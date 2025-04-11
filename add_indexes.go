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

	log.Println("Database indexes applied successfully")
} 