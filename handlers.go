package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func createTag(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tag Tag
		if err := c.ShouldBindJSON(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get organizationId from path
		orgID, err := strconv.ParseUint(c.Param("organizationId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}
		tag.OrganizationID = uint(orgID)

		if err := db.Create(&tag).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, tag)
	}
}

func createFinancialRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var record FinancialRecord
		if err := c.ShouldBindJSON(&record); err != nil {
			fmt.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get organizationId from path
		orgID, err := strconv.ParseUint(c.Param("organizationId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}
		record.OrganizationID = uint(orgID)

		// Validate direction
		if record.Direction != "IN" && record.Direction != "OUT" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Direction must be either 'IN' or 'OUT'"})
			return
		}

		// Validate amount
		if record.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than or equal to zero"})
			return
		}

		if err := db.Create(&record).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, record)
	}
}

func createFinancialRecordsBulk(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var records []FinancialRecord
		if err := c.ShouldBindJSON(&records); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get organizationId from path
		orgID, err := strconv.ParseUint(c.Param("organizationId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}

		// Validate and set organization ID for all records
		for i := range records {
			records[i].OrganizationID = uint(orgID)

			// Validate direction
			if records[i].Direction != "IN" && records[i].Direction != "OUT" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Direction must be either 'IN' or 'OUT'"})
				return
			}

			// Validate amount
			if records[i].Amount < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than or equal to zero"})
				return
			}
		}

		// Create all records in a single transaction
		if err := db.Create(&records).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, records)
	}
}

func listFinancialRecords(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID, err := strconv.ParseUint(c.Param("organizationId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}

		// Parse pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		// Ensure page and pageSize are positive
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 20
		}

		// Calculate offset
		offset := (page - 1) * pageSize

		query := db.Where("organization_id = ?", orgID)

		// Handle tag filtering
		if tagIDs := c.Query("tags"); tagIDs != "" {
			tagIDList := strings.Split(tagIDs, ",")
			query = query.Joins("JOIN financial_record_tags ON financial_record_tags.financial_record_id = financial_records.id").
				Where("financial_record_tags.tag_id IN ?", tagIDList)
		}

		// Get total count for pagination
		var total int64
		if err := query.Model(&FinancialRecord{}).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var records []FinancialRecord
		if err := query.Preload("Tags").
			Offset(offset).
			Limit(pageSize).
			Find(&records).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Calculate total pages
		totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

		c.JSON(http.StatusOK, gin.H{
			"data": records,
			"pagination": gin.H{
				"current_page": page,
				"page_size":    pageSize,
				"total_items":   total,
				"total_pages":  totalPages,
			},
		})
	}
}

func getCashFlowReport(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID, err := strconv.ParseUint(c.Param("organizationId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}

		// Calculate date range (last 2 years)
		now := time.Now()
		twoYearsAgo := now.AddDate(-2, 0, 0)

		// Use raw SQL to aggregate data in the database
		var monthlyData []MonthlyCashFlow
		if err := db.Raw(`
			SELECT 
				EXTRACT(YEAR FROM due_date)::integer as year,
				EXTRACT(MONTH FROM due_date)::integer as month,
				SUM(CASE WHEN direction = 'IN' THEN amount ELSE 0 END) as in,
				SUM(CASE WHEN direction = 'OUT' THEN amount ELSE 0 END) as out
			FROM financial_records
			WHERE organization_id = ? AND due_date >= ?
			GROUP BY EXTRACT(YEAR FROM due_date), EXTRACT(MONTH FROM due_date)
			ORDER BY year, month
		`, orgID, twoYearsAgo).Scan(&monthlyData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Map the database results to our response structure
		report := CashFlowReport{
			MonthlyData: monthlyData,
		}

		c.JSON(http.StatusOK, report)
	}
}

func listTags(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID, err := strconv.ParseUint(c.Param("organizationId"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			return
		}

		// Parse pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		// Ensure page and pageSize are positive
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 20
		}

		// Calculate offset
		offset := (page - 1) * pageSize

		// Get total count for pagination
		var total int64
		if err := db.Model(&Tag{}).Where("organization_id = ?", orgID).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var tags []Tag
		if err := db.Where("organization_id = ?", orgID).
			Offset(offset).
			Limit(pageSize).
			Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Calculate total pages
		totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

		c.JSON(http.StatusOK, gin.H{
			"data": tags,
			"pagination": gin.H{
				"current_page": page,
				"page_size":    pageSize,
				"total_items":  total,
				"total_pages": totalPages,
			},
		})
	}
}