package waitlisthandler

import (
	"net/http"
	"strconv"
	"time"

	waitlistdto "github.com/demola234/defifundr/internal/features/waitlist/dto"
	waitlistport "github.com/demola234/defifundr/internal/features/waitlist/port"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	appErrors "github.com/demola234/defifundr/pkg/apperrors"
	"github.com/demola234/defifundr/pkg/metrics"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service waitlistport.Service
	logger  logging.Logger
}

func New(service waitlistport.Service, logger logging.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) JoinWaitlist(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("waitlist-handler").Start(ctx.Request.Context(), "JoinWaitlist")
	defer span.End()
	ctx.Request = ctx.Request.WithContext(spanCtx)

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing join waitlist request")

	var req waitlistdto.JoinRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, waitlistdto.ErrorResponse{Message: appErrors.ErrInvalidRequest.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, waitlistdto.ErrorResponse{Message: appErrors.ErrInvalidRequest.Error()})
		return
	}

	entry, err := h.service.JoinWaitlist(ctx.Request.Context(), req.Email, req.FullName, req.ReferralSource)
	if err != nil {
		if appErrors.IsAppError(err) {
			appErr := err.(*appErrors.AppError)
			if appErr.ErrorType == appErrors.ErrorTypeConflict {
				ctx.JSON(http.StatusConflict, waitlistdto.ErrorResponse{Message: appErr.Error()})
				return
			}
			ctx.JSON(http.StatusBadRequest, waitlistdto.ErrorResponse{Message: appErr.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, waitlistdto.ErrorResponse{Message: appErrors.ErrInternalServer.Error()})
		return
	}

	position, err := h.service.GetWaitlistPosition(ctx.Request.Context(), entry.ID)
	if err != nil {
		reqLogger.Warn("Failed to get waitlist position", map[string]any{"error": err.Error()})
	}

	metrics.WaitlistSignupsTotal.Inc()
	ctx.JSON(http.StatusCreated, waitlistdto.SuccessResponse{
		Success: true,
		Message: "Successfully joined waitlist",
		Data: waitlistdto.EntryResponse{
			ID:             entry.ID,
			Email:          entry.Email,
			FullName:       entry.FullName,
			ReferralCode:   entry.ReferralCode,
			ReferralSource: entry.ReferralSource,
			Status:         entry.Status,
			Position:       position,
			SignupDate:     entry.SignupDate,
		},
	})
}

func (h *Handler) ListWaitlist(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("waitlist-handler").Start(ctx.Request.Context(), "ListWaitlist")
	defer span.End()
	ctx.Request = ctx.Request.WithContext(spanCtx)

	role, exists := ctx.Get("user_role")
	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, waitlistdto.ErrorResponse{Message: "Access denied"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filters := make(map[string]string)
	if s := ctx.Query("status"); s != "" {
		filters["status"] = s
	}
	if src := ctx.Query("source"); src != "" {
		filters["referral_source"] = src
	}
	if ord := ctx.Query("order"); ord != "" {
		filters["order"] = ord
	}

	entries, total, err := h.service.ListWaitlist(ctx.Request.Context(), page, pageSize, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, waitlistdto.ErrorResponse{Message: "Failed to retrieve waitlist entries"})
		return
	}

	items := make([]waitlistdto.EntryResponse, len(entries))
	for i, e := range entries {
		items[i] = waitlistdto.EntryResponse{
			ID:             e.ID,
			Email:          e.Email,
			FullName:       e.FullName,
			ReferralCode:   e.ReferralCode,
			ReferralSource: e.ReferralSource,
			Status:         e.Status,
			SignupDate:     e.SignupDate,
			InvitedDate:    e.InvitedDate,
		}
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, waitlistdto.PageResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Items:      items,
	})
}

func (h *Handler) GetWaitlistStats(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("waitlist-handler").Start(ctx.Request.Context(), "GetWaitlistStats")
	defer span.End()
	ctx.Request = ctx.Request.WithContext(spanCtx)

	role, exists := ctx.Get("user_role")
	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, waitlistdto.ErrorResponse{Message: "Access denied"})
		return
	}

	stats, err := h.service.GetWaitlistStats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, waitlistdto.ErrorResponse{Message: "Failed to retrieve waitlist statistics"})
		return
	}

	ctx.JSON(http.StatusOK, waitlistdto.SuccessResponse{
		Success: true,
		Message: "Waitlist statistics retrieved",
		Data:    stats,
	})
}

func (h *Handler) ExportWaitlist(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("waitlist-handler").Start(ctx.Request.Context(), "ExportWaitlist")
	defer span.End()
	ctx.Request = ctx.Request.WithContext(spanCtx)

	role, exists := ctx.Get("user_role")
	if !exists || role != "admin" {
		ctx.JSON(http.StatusForbidden, waitlistdto.ErrorResponse{Message: "Access denied"})
		return
	}

	csvData, err := h.service.ExportWaitlist(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, waitlistdto.ErrorResponse{Message: "Failed to export waitlist data"})
		return
	}

	filename := "waitlist-export-" + time.Now().Format("2006-01-02") + ".csv"
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Length", strconv.Itoa(len(csvData)))
	ctx.Data(http.StatusOK, "text/csv", csvData)
}
