package server

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handlerSmartInfo(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	res, err := s.smart.GetNVMeSmart(id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(res)
}
