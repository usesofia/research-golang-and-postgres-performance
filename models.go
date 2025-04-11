package main

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	OrganizationID uint   `json:"organizationId" gorm:"not null"`
	Name           string `json:"name" gorm:"not null"`
}

type FinancialRecord struct {
	gorm.Model
	OrganizationID uint      `json:"organizationId" gorm:"not null"`
	Direction      string    `json:"direction" gorm:"not null"` // "IN" or "OUT"
	Amount         float64   `json:"amount" gorm:"not null"`
	Tags           []Tag     `json:"tags" gorm:"many2many:financial_record_tags;"`
	DueDate        time.Time `json:"dueDate" gorm:"not null"`
}

type CashFlowReport struct {
	MonthlyData []MonthlyCashFlow `json:"monthlyData"`
}

type MonthlyCashFlow struct {
	Year  int     `json:"year"`
	Month int     `json:"month"`
	In    float64 `json:"in"`
	Out   float64 `json:"out"`
}
