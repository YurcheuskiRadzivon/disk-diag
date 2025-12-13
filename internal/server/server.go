package server

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 40 * time.Second
	_defaultWriteTimeout    = 40 * time.Second
	_defaultShutdownTimeout = 40 * time.Second
)

type Server struct {
	App     *fiber.App
	notify  chan error
	address string
}

func New(port string) *Server {
	s := &Server{
		App:     nil,
		notify:  make(chan error, 1),
		address: port,
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    30 * 1024 * 1024,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	s.App = app

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.App.Listen(s.address)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.App.ShutdownWithTimeout(_defaultShutdownTimeout)
}
