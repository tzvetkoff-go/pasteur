package webserver

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	"github.com/tzvetkoff-go/errors"
	"github.com/tzvetkoff-go/logger"

	"github.com/tzvetkoff-go/pasteur/pkg/codemirror"
	"github.com/tzvetkoff-go/pasteur/pkg/config"
	"github.com/tzvetkoff-go/pasteur/pkg/db"
	"github.com/tzvetkoff-go/pasteur/pkg/hasher"
	"github.com/tzvetkoff-go/pasteur/pkg/httplib"
	"github.com/tzvetkoff-go/pasteur/pkg/model"
)

// StaticFSRoot ...
//go:embed static
var StaticFSRoot embed.FS

// TemplatesFSRoot ...
//go:embed templates
var TemplatesFSRoot embed.FS

// WebServer ...
type WebServer struct {
	ListenAddress string
	TLSCert       string
	TLSKey        string

	App             *fiber.App
	Views           *html.Engine
	DefaultTemplate *template.Template

	Hasher *hasher.Hasher `inject:"Hasher"`
	DB     db.DB          `inject:"DB"`
}

// Variables ...
type Variables struct {
	Paste     *model.Paste
	Languages []codemirror.Mode
}

// New ...
func New(wsConfig *config.WebServer) (*WebServer, error) {
	result := &WebServer{
		ListenAddress: wsConfig.ListenAddress,
		TLSCert:       wsConfig.TLSCert,
		TLSKey:        wsConfig.TLSKey,
	}

	var err error

	var staticFS fs.FS
	if wsConfig.StaticPath == "embedded" {
		staticFS, err = fs.Sub(StaticFSRoot, "static")
		if err != nil {
			return nil, errors.Propagate(err, "cannot access embedded static/ subdirectory")
		}
	} else {
		staticFS = os.DirFS(wsConfig.StaticPath)
	}

	var viewsEngine *html.Engine

	if wsConfig.TemplatesPath == "embedded" {
		templatesFS, err := fs.Sub(TemplatesFSRoot, "templates")
		if err != nil {
			return nil, errors.Propagate(err, "cannot access embedded templates/ directory")
		}

		viewsEngine = html.NewFileSystem(http.FS(templatesFS), ".html")
	} else {
		viewsEngine = html.New(wsConfig.TemplatesPath, ".html")
	}
	result.Views = viewsEngine

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 viewsEngine,
		ErrorHandler:          httplib.ErrorHandler,
		ProxyHeader:           wsConfig.ProxyHeader,
	})

	app.Use(httplib.RequestId())
	app.Use(httplib.RequestLogger())
	app.Use(httplib.ErrorRecoverer())

	app.Get("/", httplib.Timeout(result.Index, 10*time.Second))
	app.Get("/p/:id.txt", httplib.Timeout(result.ShowRaw, 10*time.Second))
	app.Get("/p/:id", httplib.Timeout(result.Show, 10*time.Second))
	app.Post("/", httplib.Timeout(result.Create, 10*time.Second))

	app.Get("/*", httplib.Timeout(filesystem.New(filesystem.Config{
		Root: http.FS(staticFS),
	}), 10*time.Second))

	app.Use(httplib.NotFoundHandler)

	result.App = app

	return result, nil
}

// Index ...
func (ws *WebServer) Index(c *fiber.Ctx) error {
	return c.Render("default/page", fiber.Map{
		"Paste":     model.NewPaste(),
		"Languages": codemirror.Modes,
	}, "layout/layout")
}

// Show ...
func (ws *WebServer) Show(c *fiber.Ctx) error {
	idString := c.Params("id")
	id, err := ws.Hasher.Decode(idString)
	if err != nil {
		// Mask error - we don't want to expose hasher.Decode errors
		logger.Error("%s", err)
		return httplib.NotFoundHandler(c)
	}

	paste, err := ws.DB.RetrievePasteByID(id)
	if err != nil {
		return err
	}

	return c.Render("default/page", fiber.Map{
		"Paste":     paste,
		"Languages": codemirror.Modes,
	}, "layout/layout")
}

// ShowRaw ...
func (ws *WebServer) ShowRaw(c *fiber.Ctx) error {
	idString := c.Params("id")
	id, err := ws.Hasher.Decode(idString)
	if err != nil {
		// Mask error - we don't want to expose hasher.Decode errors
		logger.Error("%s", err)
		return httplib.NotFoundHandler(c)
	}

	paste, err := ws.DB.RetrievePasteByID(id)
	if err != nil {
		return err
	}

	c.Response().Header.Set("X-Filename", paste.Filename)
	c.Response().Header.Set("X-MIME-Type", paste.MimeType)
	return c.SendString(paste.Content)
}

// Create ...
func (ws *WebServer) Create(c *fiber.Ctx) error {
	paste := model.NewPaste()

	contentType := string(c.Request().Header.Peek("ContentType"))
	if strings.Contains(contentType, "json") {
		err := json.Unmarshal(c.Request().Body(), &paste)
		if err != nil {
			// Mask error - we don't want to expose json.Unmarshal errors
			logger.Error("%s", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}
	} else {
		paste.IndentStyle = c.FormValue("indent-style")
		paste.IndentSize = c.FormValue("indent-size")
		paste.MimeType = c.FormValue("mime-type")
		paste.Filename = c.FormValue("filename")
		paste.Content = c.FormValue("content")
	}

	err := paste.Validate()
	if err != nil {
		// Send proper error status
		c.Status(fiber.StatusUnprocessableEntity)
		return c.SendString(fmt.Sprintf("%#s", err))
	}

	paste, err = ws.DB.CreatePaste(paste)
	if err != nil {
		return err
	}

	if strings.Contains(contentType, "json") {
		return c.JSON(paste)
	}

	hash, err := ws.Hasher.Encode(paste.ID)
	if err != nil {
		return err
	}

	return c.Redirect(fmt.Sprintf("/p/%s", hash), fiber.StatusFound)
}

// Serve ...
func (ws *WebServer) Serve() error {
	if ws.TLSCert != "" && ws.TLSKey != "" {
		logger.Info("Starting TLS web server at %s", ws.ListenAddress)
		return ws.App.ListenTLS(ws.ListenAddress, ws.TLSCert, ws.TLSKey)
	}

	logger.Info("Starting web server at %s", ws.ListenAddress)
	return ws.App.Listen(ws.ListenAddress)
}
