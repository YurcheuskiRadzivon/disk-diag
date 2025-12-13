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
	parts, _ := s.base.GetExtendedPartitions(id)

	return ctx.Status(http.StatusOK).JSON(parts)
}
