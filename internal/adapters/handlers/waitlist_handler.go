// internal/adapters/handlers/waitlist_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/adapters/dto/request"
	"github.com/demola234/defifundr/internal/adapters/dto/response"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/app_errors"
	"github.com/gin-gonic/gin"
)

type WaitlistHandler struct {
	waitlistService ports.WaitlistService
	logger          logging.Logger
}

// NewWaitlistHandler creates a new waitlist handler
func NewWaitlistHandler(waitlistService ports.WaitlistService, logger logging.Logger) *WaitlistHandler {
	return &WaitlistHandler{
		waitlistService: waitlistService,
		logger:          logger,
	}
}

// JoinWaitlist godoc
// @Summary Join the waitlist
// @Description Register for early access to the platform
// @Tags waitlist
// @Accept json
// @Produce json
// @Param join body request.WaitlistJoinRequest true "Waitlist join data"
// @Success 201 {object} response.WaitlistEntryResponse "Successfully joined waitlist"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 409 {object} response.ErrorResponse "Email already on waitlist"
// @Failure 429 {object} response.ErrorResponse "Too many requests"
// @Router /waitlist [post]
func (h *WaitlistHandler) JoinWaitlist(ctx *gin.Context) {
	// Extract request co-relation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing join waitlist request")

	var req request.WaitlistJoinRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   app_errors.ErrInvalidRequest.Error(),
			Details: err.Error(),
		})
		return
	}

	// Validate request data
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   app_errors.ErrInvalidRequest.Error(),
			Details: err.Error(),
		})
		return
	}

	// Join the waitlist
	entry, err := h.waitlistService.JoinWaitlist(ctx, req.Email, req.FullName, req.ReferralSource)
	if err != nil {
		errResponse := response.ErrorResponse{
			Error: app_errors.ErrInternalServer.Error(),
		}

		if app_errors.IsAppError(err) {
			appErr := err.(*app_errors.AppError)
			errResponse.Error = appErr.Error()

			if appErr.ErrorType == app_errors.ErrorTypeConflict {
				ctx.JSON(http.StatusConflict, errResponse)
				return
			}

			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	// Get waitlist position
	position, err := h.waitlistService.GetWaitlistPosition(ctx, entry.ID)
	if err != nil {
		// Non-critical error, continue with response
		reqLogger.Warn("Failed to get waitlist position", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Create response
	entryResponse := response.WaitlistEntryResponse{
		ID:             entry.ID,
		Email:          entry.Email,
		FullName:       entry.FullName,
		ReferralCode:   entry.ReferralCode,
		ReferralSource: entry.ReferralSource,
		Status:         entry.Status,
		Position:       position,
		SignupDate:     entry.SignupDate,
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		Message: "Successfully joined waitlist",
		Data:    entryResponse,
	})
}


// ListWaitlist godoc
// @Summary List waitlist entries
// @Description List waitlist entries with pagination and filtering
// @Tags waitlist
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param status query string false "Filter by status (waiting, invited, registered)"
// @Param source query string false "Filter by referral source"
// @Param order query string false "Order by (signup_date_asc, signup_date_desc)"
// @Success 200 {object} response.PageResponse "Paginated list of waitlist entries"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Security BearerAuth
// @Router /admin/waitlist [get]
func (h *WaitlistHandler) ListWaitlist(ctx *gin.Context) {
	// Check if user is admin (assuming auth middleware has set role)
	role, exists := ctx.Get("user_role")
	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{
			Error: "Access denied",
		})
		return
	}

	// Parse pagination parameters
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Parse filters
	status := ctx.Query("status")
	source := ctx.Query("source")
	order := ctx.Query("order")

	filters := make(map[string]string)
	if status != "" {
		filters["status"] = status
	}
	if source != "" {
		filters["referral_source"] = source
	}
	if order != "" {
		filters["order"] = order
	}

	// Get waitlist entries
	entries, total, err := h.waitlistService.ListWaitlist(ctx, page, pageSize, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to retrieve waitlist entries",
		})
		return
	}

	// Map to response objects
	responseEntries := make([]response.WaitlistEntryResponse, len(entries))
	for i, entry := range entries {
		responseEntries[i] = response.WaitlistEntryResponse{
			ID:             entry.ID,
			Email:          entry.Email,
			FullName:       entry.FullName,
			ReferralCode:   entry.ReferralCode,
			ReferralSource: entry.ReferralSource,
			Status:         entry.Status,
			SignupDate:     entry.SignupDate,
		}

		if entry.InvitedDate != nil {
			responseEntries[i].InvitedDate = entry.InvitedDate
		}
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// Create paginated response
	ctx.JSON(http.StatusOK, response.PageResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Items:      responseEntries,
	})
}

// GetWaitlistStats godoc
// @Summary Get waitlist statistics
// @Description Get statistics about the waitlist
// @Tags waitlist
// @Accept json
// @Produce json
// @Success 200 {object} response.WaitlistStatsResponse "Waitlist statistics"
// @Security BearerAuth
// @Router /admin/waitlist/stats [get]
func (h *WaitlistHandler) GetWaitlistStats(ctx *gin.Context) {
	// Check if user is admin (assuming auth middleware has set role)
	role, exists := ctx.Get("user_role")
	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{
			Error: "Access denied",
		})
		return
	}

	// Get statistics
	stats, err := h.waitlistService.GetWaitlistStats(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to retrieve waitlist statistics",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Waitlist statistics retrieved",
		Data:    stats,
	})
}

// ExportWaitlist godoc
// @Summary Export waitlist data
// @Description Export all waitlist data as CSV
// @Tags waitlist
// @Accept json
// @Produce text/csv
// @Success 200 {file} blob "CSV file"
// @Security BearerAuth
// @Router /admin/waitlist/export [get]
func (h *WaitlistHandler) ExportWaitlist(ctx *gin.Context) {
	// Check if user is admin (assuming auth middleware has set role)
	role, exists := ctx.Get("user_role")
	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{
			Error: "Access denied",
		})
		return
	}

	// Export waitlist data
	csvData, err := h.waitlistService.ExportWaitlist(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to export waitlist data",
		})
		return
	}

	// Set response headers
	currentTime := time.Now().Format("2006-01-02")
	filename := "waitlist-export-" + currentTime + ".csv"
	
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Length", strconv.Itoa(len(csvData)))
	
	ctx.Data(http.StatusOK, "text/csv", csvData)
}