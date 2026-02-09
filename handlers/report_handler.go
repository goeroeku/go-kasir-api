package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"kasir-api/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetReportToday godoc
// @Summary Get sales report for today
// @Description Get total revenue, total transactions, and best seller for today
// @Tags Reports
// @Produce json
// @Success 200 {object} models.SalesReport
// @Router /reports/today [get]
func (h *ReportHandler) GetReportToday(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetReportToday()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetReportCustom godoc
// @Summary Get sales report with custom date range
// @Description Get sales report filtered by start_date and end_date (YYYY-MM-DD)
// @Tags Reports
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} models.SalesReport
// @Router /reports [get]
func (h *ReportHandler) GetReportCustom(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		http.Error(w, "Invalid start_date format (use YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		http.Error(w, "Invalid end_date format (use YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Adjust end date to include the whole day
	endDate = endDate.Add(24 * time.Hour).Add(-1 * time.Nanosecond)

	report, err := h.service.GetReportByDateRange(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Handler routes requests to appropriate method handlers
func (h *ReportHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if it's /reports/today or /reports with query params
	// Simple routing logic for now, relies on main.go routing to distinguish path
	// But since main.go will likely route /reports/today to a separate handlerfunc or same handler struct
	// Let's assume this Handler method handles /reports generic which might be today or custom
	// Actually, better to split in main.go
}
