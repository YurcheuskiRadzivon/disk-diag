package server

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handlerDiagnosticManualInfo(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	smartInfo, err := s.smart.GetNVMeSmart(id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	res, err := s.diagnostic.AnalyzeManual(smartInfo)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(res)
}

func (s *Server) handlerDiagnosticGeminiInfo(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	smartInfo, err := s.smart.GetNVMeSmart(id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	res, err := s.diagnostic.AnalyzeAI(smartInfo)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(res)
}
