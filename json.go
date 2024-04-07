package www

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const ContentType = "application/json; charset=utf-8"

func Json(errNotFound error) fiber.Handler {
	handleResponse := func(c *fiber.Ctx) error {
		err := c.Next()

		if err == nil {
			return nil
		}

		var fe *fiber.Error
		if errors.As(err, &fe) {
			return respondError(c, fe.Code, fe.Message)
		}

		if errors.Is(err, errNotFound) {
			ErrNotFound(c)
		}

		log.Err(err).
			Msg("un handled error while processing request")

		return ErrInternal(c, err)
	}

	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		c.Set(fiber.HeaderContentType, ContentType)

		switch c.Method() {

		case fiber.MethodGet, fiber.MethodDelete:
			return handleResponse(c)

		case fiber.MethodPost, fiber.MethodPut:
			if !c.Is("json") {
				return respondError(c, fiber.StatusUnsupportedMediaType, "unsupported media type")
			}
			return handleResponse(c)

		case fiber.MethodOptions:
			return c.SendStatus(fiber.StatusNoContent)

		default:
			return ErrMethodNotAllowed(c)

		}
	}
}
