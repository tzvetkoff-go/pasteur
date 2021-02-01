package httplib

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/tzvetkoff-go/logger"
)

// NotFoundHandler ...
func NotFoundHandler(c *fiber.Ctx) error {
	return c.SendStatus(404)
}

// ErrorHandler ...
func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == sql.ErrNoRows {
		return NotFoundHandler(c)
	}

	if err, ok := err.(error); ok && os.IsNotExist(err) {
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
