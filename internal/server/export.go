package server

import (
	"net/http"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handlerExport(ctx *fiber.Ctx) error {
	var req models.ExportRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if req.Data == nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Data is required")
	}

	var (
		bytes       []byte
		contentType string
		err         error
		filename    string
	)

	switch req.Format {
	case "json":
		bytes, contentType, err = s.export.ExportJSON(req.Data)
		filename = "export.json"
	case "txt":
		bytes, contentType, err = s.export.ExportTXT(req.Data, req.Title)
		filename = "export.txt"
	case "pdf":
		bytes, contentType, err = s.export.ExportPDF(req.Data, req.Title)
		filename = "export.pdf"
	default:
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid format. Use: json, txt, pdf")
	}

	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.Set("Content-Type", contentType)
	ctx.Set("Content-Disposition", "attachment; filename="+filename)
	return ctx.Send(bytes)
}
