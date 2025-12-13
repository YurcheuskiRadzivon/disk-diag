package server

import (
	"encoding/json"
	"time"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/base"
	"github.com/gofiber/fiber/v2"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 40 * time.Second
	_defaultWriteTimeout    = 40 * time.Second
	_defaultShutdownTimeout = 40 * time.Second
)

type Server struct {
	app     *fiber.App
	notify  chan error
	address string
	base    base.Service
}

type Error struct {
	Message string `json:"message" example:"message"`
}

func New(port string, base base.Service) *Server {
	s := &Server{
		app:     nil,
		notify:  make(chan error, 1),
		address: port,
		base:    base,
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    30 * 1024 * 1024,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	s.app = app

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.app.Listen(s.address)
		close(s.notify)
	}()
}

func (s *Server) RegisterRoutes() {
	base := s.app.Group("/base")
	{
		base.Get("/disks", s.handleDisks)
		base.Get("/cdiskinfo/:id", s.handleCDiskInfo)
		base.Get("/partitions/:id", s.handlePartitions)

	}
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.app.ShutdownWithTimeout(_defaultShutdownTimeout)
}

func ErrorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(Error{Message: msg})
}
