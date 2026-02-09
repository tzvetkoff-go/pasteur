package httplib

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/tzvetkoff-go/logger"
)

// NotFoundHandler ...
func NotFoundHandler(c fiber.Ctx) error {
	return c.SendStatus(404)
}

// ErrorHandler ...
func ErrorHandler(c fiber.Ctx, err error) error {
	if err == sql.ErrNoRows {
		return NotFoundHandler(c)
	}

	if os.IsNotExist(err) {
		return NotFoundHandler(c)
	}

	if e, ok := err.(*fiber.Error); ok {
		if e.Code != 404 {
			logger.Error("%s", err)
		}

		return c.SendStatus(e.Code)
	}

	logger.Error("%s", err)
	return c.SendStatus(fiber.StatusInternalServerError)
}

// Redirect ...
func Redirect(path string, code int) fiber.Handler {
	fn := func(c fiber.Ctx) error {
		return c.Redirect().Status(code).To(path)
	}

	return fn
}
