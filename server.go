package www

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
)

type (
	ServerConfig struct {
		Port      int
		Timeout   time.Duration
		BodyLimit int
		Logs      bool
		Origins   []string
	}

	Route func(r fiber.Router)

	Option func(*ServerConfig)
)

// Serve starts an instance of the http server.
func Serve(route Route, config ServerConfig) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		CaseSensitive:           true,
		StreamRequestBody:       true,
		ReadTimeout:             config.Timeout,
		WriteTimeout:            config.Timeout,
		BodyLimit:               config.BodyLimit,
		EnableTrustedProxyCheck: true,
	})
	if config.Logs {
		app.Use(logger.New())
	}

	if len(config.Origins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins: strings.Join(config.Origins, ", "),
		}))
	}

	group := app.Group("/")
	route(group)
	addr := fmt.Sprintf(":%d", config.Port)
	log.Info().Str("addr", addr).Msg("starting http server")
	return app.Listen(addr)
}
