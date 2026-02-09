package httplib

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/tzvetkoff-go/logger"
)

// ErrorRecoverer ...
func ErrorRecoverer() func(fiber.Ctx) error {
	fn := func(c fiber.Ctx) error {
		var err error

		defer func() {
			if r := recover(); r != nil {
				if r == sql.ErrNoRows {
					NotFoundHandler(c)
					return
				}

				var ok bool
				if err, ok = r.(error); ok && os.IsNotExist(err) {
					NotFoundHandler(c)
					return
				}

				ErrorHandler(c, r.(error))
			}
		}()

		if err != nil {
			return err
		}

		return c.Next()
	}

	return fn
}

// RequestId ...
func RequestId() func(fiber.Ctx) error {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	var buf [12]byte
	var b64 string

	for len(b64) < 8 {
		_, _ = rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix := fmt.Sprintf("%s/%s", hostname, b64[0:8])
	var seq uint64

	fn := func(c fiber.Ctx) error {
		requestId := string(c.Request().Header.Peek("X-Request-Id"))
		if requestId == "" {
			atomic.AddUint64(&seq, 1)
			requestId = fmt.Sprintf("%s/%08x", prefix, seq)
		}

		c.Set("X-Request-Id", requestId)

		return c.Next()
	}

	return fn
}

// RequestLogger ...
func RequestLogger() func(fiber.Ctx) error {
	// debugFlag := (logger.GetLevel() & logger.LOG_DEBUG) == logger.LOG_DEBUG

	fn := func(c fiber.Ctx) error {
		fields := logger.Fields{}

		fields["scheme"] = "http"
		if c.RequestCtx().IsTLS() {
			fields["scheme"] = "https"
		}

		fields["host"] = c.Hostname()
		fields["method"] = string(c.Request().Header.Method())

		requestPath := string(c.Request().URI().Path())
		queryString := string(c.Request().URI().QueryString())
		if queryString != "" {
			requestPath += "?" + queryString
		}
		fields["path"] = requestPath

		fields["ip"] = c.IP()
		fields["ips"] = c.IPs()

		t := time.Now()
		err := c.Next()
		if err != nil {
			if c.App().Config().ErrorHandler(c, err) != nil {
				c.SendStatus(500)
			}
		}
		fields["duration"] = time.Since(t)

		fields["status"] = c.Response().StatusCode()
		fields["request_size"] = len(c.Request().Body())
		fields["response_size"] = len(c.Response().Body())

		// if debugFlag && !strings.HasPrefix(requestPath, "/assets/") && requestPath != "/favicon.ico" {
		// 	fields["request_body"] = strings.TrimSpace(string(c.Request().Body()))
		// 	fields["response_body"] = strings.TrimSpace(string(c.Response().Body()))
		// }

		if c.Response().StatusCode() >= 400 {
			logger.Error("", fields)
		} else {
			logger.Info("", fields)
		}

		return nil
	}

	return fn
}

// Timeout ...
func Timeout(handler fiber.Handler, timeout time.Duration) fiber.Handler {
	fn := func(ctx fiber.Ctx) error {
		var err error
		ch := make(chan struct{}, 1)

		go func() {
			defer func() {
				_ = recover()
			}()

			err = handler(ctx)
			ch <- struct{}{}
		}()

		select {
		case <-ch:
		case <-time.After(timeout):
			return fiber.ErrRequestTimeout
		}

		return err
	}

	return fn
}
