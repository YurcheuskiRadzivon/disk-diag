package server

import (
	"encoding/json"
	"time"

	"net/http"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/base"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/benchmark"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/diagnostic"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/export"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/smart"
	"github.com/YurcheuskiRadzivon/disk-diag/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 40 * time.Second
	_defaultWriteTimeout    = 40 * time.Second
	_defaultShutdownTimeout = 40 * time.Second
)

type Server struct {
	app        *fiber.App
	notify     chan error
	address    string
	base       base.Service
	smart      smart.Service
	benchmark  benchmark.Service
	diagnostic diagnostic.Service
	export     export.Service
}

type Error struct {
	Message string `json:"message" example:"message"`
}

func New(port string, base base.Service, smart smart.Service, benchmark benchmark.Service, diagnostic diagnostic.Service) *Server {
	s := &Server{
		app:        nil,
		notify:     make(chan error, 1),
		address:    port,
		base:       base,
		smart:      smart,
		benchmark:  benchmark,
		diagnostic: diagnostic,
		export:     export.New(),
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
	ui := s.app.Group("/ui")
	{
		ui.Use("/", filesystem.New(filesystem.Config{
			Root:       http.FS(web.Assets),
			PathPrefix: "",
			Browse:     true,
		}))
	}

	base := s.app.Group("/base")
	{
		base.Get("/disks", s.handleDisks)
		base.Get("/cdiskinfo/:id", s.handleCDiskInfo)
		base.Get("/partitions/:id", s.handlePartitions)

	}

	smart := s.app.Group("/smart")
	{
		smart.Get("/:id", s.handlerSmartInfo)

	}

	benchmark := s.app.Group("/benchmark")
	{
		benchmark.Get("/write/:id", s.handlerBenchmarkWriteInfo)
		benchmark.Get("/read/:id", s.handlerBenchmarkReadInfo)
		benchmark.Get("/iops/:id", s.handlerBenchmarkIOPSInfo)

	}

	diagnostic := s.app.Group("/diagnostic")
	{
		diagnostic.Get("/gemini/:id", s.handlerDiagnosticGeminiInfo)
		diagnostic.Get("/diagnostic/:id", s.handlerDiagnosticManualInfo)
	}

	exportGroup := s.app.Group("/export")
	{
		exportGroup.Post("/", s.handlerExport)
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
