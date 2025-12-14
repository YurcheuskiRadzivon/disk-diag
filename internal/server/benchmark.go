package server

import (
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handlerBenchmarkWriteInfo(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	cfg := models.TestConfig{}
	res, err := s.benchmark.WriteTest(id, &cfg)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (s *Server) handlerBenchmarkReadInfo(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	cfg := models.TestConfig{
		Retry: true,
	}
	_, err := s.benchmark.WriteTest(id, &cfg)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	res, err := s.benchmark.ReadTest(id, &cfg)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (s *Server) handlerBenchmarkIOPSInfo(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	cfg := models.TestConfig{
		Retry: true,
	}
	_, err := s.benchmark.WriteTest(id, &cfg)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	res, err := s.benchmark.IOPSTest(id, &cfg)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(res)
}
