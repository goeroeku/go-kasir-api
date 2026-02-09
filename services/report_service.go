package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReportToday() (*models.SalesReport, error) {
	now := time.Now()
	// Start of day
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// End of day
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	return s.repo.GetSalesReport(startDate, endDate)
}

func (s *ReportService) GetReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	return s.repo.GetSalesReport(startDate, endDate)
}
