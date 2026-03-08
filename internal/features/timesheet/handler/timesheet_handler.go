package timesheethandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) NotImplemented(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"success": false, "message": "coming soon"})
}
