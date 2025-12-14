package server

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleDisks(ctx *fiber.Ctx) error {
	disks, err := s.base.GetPhysicalDisks()
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(disks)
}

func (s *Server) handlePartitions(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	parts, err := s.base.GetExtendedPartitions(id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(parts)
}

func (s *Server) handleCDiskInfo(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	cDiskInfo, err := s.base.GetCDiskInfo(id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(cDiskInfo)
}
