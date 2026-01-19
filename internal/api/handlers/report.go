package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type ReportHandler struct {
	reports *service.ReportService
}

func NewReportHandler(reports *service.ReportService) *ReportHandler {
	return &ReportHandler{reports: reports}
}

func (h *ReportHandler) List(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)
	if page < 1 || pageSize < 1 || pageSize > 100 {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	offset := (page - 1) * pageSize
	reports, err := h.reports.List(c.Request.Context(), pageSize, offset)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(gin.H{
		"items": reports,
		"page":  page,
		"size":  pageSize,
	}))
}

func (h *ReportHandler) Get(c *gin.Context) {
	id := c.Param("id")
	report, err := h.reports.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}
	model.JSON(c, http.StatusOK, model.OK(report))
}

func (h *ReportHandler) Download(c *gin.Context) {
	id := c.Param("id")
	report, err := h.reports.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	path, ok := h.reports.ResolvePath(report.FilePath)
	if !ok {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	if _, err := os.Stat(path); err != nil {
		model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
		return
	}

	c.FileAttachment(path, report.Name)
}

func (h *ReportHandler) Preview(c *gin.Context) {
	id := c.Param("id")
	report, err := h.reports.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	if report.Type != "html" {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	path, ok := h.reports.ResolvePath(report.FilePath)
	if !ok {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	if _, err := os.Stat(path); err != nil {
		model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
		return
	}

	c.File(path)
}
