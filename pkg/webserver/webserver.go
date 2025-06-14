package webserver

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/tzvetkoff-go/errors"
	"github.com/tzvetkoff-go/logger"
	"github.com/valyala/fasthttp"

	"github.com/tzvetkoff-go/pasteur/pkg/db"
	"github.com/tzvetkoff-go/pasteur/pkg/hasher"
	"github.com/tzvetkoff-go/pasteur/pkg/httplib"
	"github.com/tzvetkoff-go/pasteur/pkg/model"
	"github.com/tzvetkoff-go/pasteur/pkg/monaco"
	"github.com/tzvetkoff-go/pasteur/pkg/stringutil"
	"github.com/tzvetkoff-go/pasteur/pkg/version"
)

// StaticFSRoot ...
//
//go:embed static
var StaticFSRoot embed.FS

// TemplatesFSRoot ...
//
//go:embed templates
var TemplatesFSRoot embed.FS

// WebServer ...
type WebServer struct {
	ListenAddress   string
	TLSCert         string
	TLSKey          string
	RelativeURLRoot string
	AbsoluteURLRoot string

	App   *fiber.App
	Views *html.Engine

	Hasher *hasher.Hasher `inject:"Hasher"`
	DB     db.DB          `inject:"DB"`
}

// APIDoc ...
// revive:disable:line-length-limit
const APIDoc = `
# Pasteur

Pasteur is a simple & stupid web pastebin.


## POST {{root}}

Creates a new paste.
Request fields:

+--------------+--------+----------+-----------+-------------------------------------------+
| Name         | Type   | Required | Default   | Description                               |
+--------------+--------+----------+-----------+-------------------------------------------+
| filename     | string | Yes      | paste.txt | Filename.                                 |
| content      | string | Yes      |           | Content.                                  |
| private      | int    | No       | 0         | Pass 1 to hide in /browse.                |
| filetype     | string | No       | automatic | Filetype.                                 |
| indent-style | string | No       | spaces    | Indent style.  Either "tabs" or "spaces". |
| indent-size  | int    | No       | tabs      | Indent size.  1..8.                       |
+--------------+--------+----------+-----------+-------------------------------------------+
| f            | file   | No       |           | Upload a file directly.                   |
+--------------+--------+----------+-----------+-------------------------------------------+

Examples:

  $ curl -XPOST {{root}} -F filename=hello-world.go -F content=$'package main\n\nfunc main() {\n\tprintln("Hello, world!")\n}\n'
  $ curl -XPOST {{root}} -F f=@hello-world.go


## GET {{root}}/browse

Lists public pastes.
Query parameters:

+----------+--------+---------+--------------------------+
| Name     | Type   | Default | Description              |
+----------+--------+---------+--------------------------+
| filetype | string |         | List pastes by filetype. |
| page     | int    | 1       | Page number.             |
| per      | int    | 20      | Page size.  1..500.      |
+----------+--------+---------+--------------------------+

Example:

  $ curl {{root}}/browse?language=go


## GET {{root}}/:id

Fetches paste content.
Example:

  $ curl {{root}}/{{id}}

`

// DefaultTimeout ...
const DefaultTimeout = 10 * time.Second

// New ...
func New(config *Config) (*WebServer, error) {
	result := &WebServer{
		ListenAddress:   config.ListenAddress,
		TLSCert:         config.TLSCert,
		TLSKey:          config.TLSKey,
		RelativeURLRoot: strings.TrimRight(config.RelativeURLRoot, "/"),
		AbsoluteURLRoot: strings.TrimRight(config.RelativeURLRoot, "/"),
	}

	if result.AbsoluteURLRoot == "" {
		result.AbsoluteURLRoot = "/"
	}

	var err error

	var staticFS fs.FS
	if config.StaticPath == "embedded" {
		staticFS, err = fs.Sub(StaticFSRoot, "static")
		if err != nil {
			return nil, errors.Propagate(err, "cannot access embedded static/ subdirectory")
		}
	} else {
		staticFS = os.DirFS(config.StaticPath)
	}

	var viewsEngine *html.Engine

	if config.TemplatesPath == "embedded" {
		templatesFS, err := fs.Sub(TemplatesFSRoot, "templates")
		if err != nil {
			return nil, errors.Propagate(err, "cannot access embedded templates/ directory")
		}

		viewsEngine = html.NewFileSystem(http.FS(templatesFS), ".html")
	} else {
		viewsEngine = html.New(config.TemplatesPath, ".html")
	}

	debugFlag := (logger.GetLevel() & logger.LOG_DEBUG) == logger.LOG_DEBUG
	if debugFlag {
		viewsEngine.Reload(true)
		viewsEngine.Debug(true)
	}

	viewsEngine.AddFunc("list", func(items ...interface{}) []interface{} {
		result := []interface{}{}
		for _, v := range items {
			result = append(result, v)
		}

		return result
	})
	viewsEngine.AddFunc("iterate", func(start int, stop int, step ...int) []int {
		result := []int{}

		ii := 1
		if len(step) > 0 {
			ii = step[0]
		}

		for i := start; i < stop; i += ii {
			result = append(result, i)
		}

		return result
	})
	viewsEngine.AddFunc("hash_encode", func(id int) string {
		encoded, _ := result.Hasher.Encode(id)
		return encoded
	})
	viewsEngine.AddFunc("hash_decode", func(encoded string) int {
		id, _ := result.Hasher.Decode(encoded)
		return id
	})
	viewsEngine.AddFunc("add_query", func(query *fasthttp.Args, name string, value interface{}) string {
		newQuery := &fasthttp.Args{}
		query.CopyTo(newQuery)

		strValue := fmt.Sprint(value)
		if strValue == "" {
			newQuery.Del(name)
		} else {
			newQuery.Set(name, strValue)
		}

		strResult := newQuery.String()
		if strResult != "" {
			return "?" + strResult
		}

		return ""
	})
	viewsEngine.AddFunc("format_filetype", func(filetype string) string {
		for _, language := range monaco.Languages {
			if language.ID == filetype {
				return language.Name
			}
		}

		return "!!" + filetype
	})
	viewsEngine.AddFunc("short_datetime", func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	})
	viewsEngine.AddFunc("full_datetime", func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05 -0700")
	})

	result.Views = viewsEngine

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 viewsEngine,
		ErrorHandler:          httplib.ErrorHandler,
		ProxyHeader:           config.ProxyHeader,
	})

	app.Use(httplib.RequestId())
	app.Use(httplib.RequestLogger())
	app.Use(httplib.ErrorRecoverer())

	if result.RelativeURLRoot != "" {
		app.Get("/", httplib.Redirect(result.RelativeURLRoot))
	}

	app.Get(result.RelativeURLRoot, httplib.Timeout(result.Index, DefaultTimeout))
	app.Get(result.RelativeURLRoot+"/browse", httplib.Timeout(result.Browse, DefaultTimeout))
	app.Get(result.RelativeURLRoot+"/:id.txt", httplib.Timeout(result.ShowRaw, DefaultTimeout))
	app.Get(result.RelativeURLRoot+"/:id", httplib.Timeout(result.Show, DefaultTimeout))
	app.Post(result.RelativeURLRoot+":/id", httplib.Timeout(result.Update, DefaultTimeout))
	app.Post(result.AbsoluteURLRoot, httplib.Timeout(result.Create, DefaultTimeout))

	fs := filesystem.New(filesystem.Config{
		Root: http.FS(staticFS),
	})

	if result.RelativeURLRoot != "" {
		app.Get(result.RelativeURLRoot+"/*", func(c *fiber.Ctx) error {
			c.Path(strings.TrimPrefix(c.Path(), result.RelativeURLRoot))
			return fs(c)
		})
	} else {
		app.Get(result.RelativeURLRoot+"/*", httplib.Timeout(fs, DefaultTimeout))
	}

	app.Use(httplib.NotFoundHandler)

	result.App = app

	return result, nil
}

// Index ...
func (ws *WebServer) Index(c *fiber.Ctx) error {
	if strings.Index(string(c.Context().UserAgent()), "curl/") == 0 {
		s := APIDoc
		s = strings.ReplaceAll(s, "{{root}}", c.BaseURL()+ws.RelativeURLRoot)
		s = strings.ReplaceAll(s, "{{id}}", ws.Hasher.EncodeAtomic(1))
		return c.SendString(s)
	}

	return c.Render("paste/page", fiber.Map{
		"RelativeURLRoot": ws.RelativeURLRoot,
		"AbsoluteURLRoot": ws.AbsoluteURLRoot,
		"ActiveMenu":      1,
		"Action":          "index",
		"Languages":       monaco.Languages,
		"Version":         version.Version,
		"Paste":           model.NewPaste(),
	}, "layout/layout")
}

// Browse ...
func (ws *WebServer) Browse(c *fiber.Ctx) error {
	conditions := map[string]interface{}{
		"private": 0,
	}

	if filetypeParam := c.Query("filetype"); filetypeParam != "" {
		conditions["filetype"] = filetypeParam
	}

	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		page = stringutil.ParseInt(pageParam)
	}

	perPage := 20
	if perPageParam := c.Query("per"); perPageParam != "" {
		perPage = stringutil.ParseInt(perPageParam)
		if perPage > 500 {
			perPage = 500
		}
	}

	order := "DESC"
	if orderParam := c.Query("order"); strings.ToUpper(orderParam) == "ASC" {
		order = "ASC"
	}

	paginatedPasteList, err := ws.DB.PaginatePastes(page, perPage, 2, order, conditions)
	if err != nil {
		return err
	}

	if strings.Index(string(c.Context().UserAgent()), "curl/") == 0 {
		t := table.NewWriter()
		s := table.StyleDefault
		s.Format.Header = text.FormatDefault
		s.Format.Footer = text.FormatDefault
		t.SetStyle(s)
		t.AppendHeader(table.Row{"Filename", "Datetime", "URL"})

		for _, paste := range paginatedPasteList.Pastes {
			t.AppendRow(table.Row{
				paste.Filename,
				paste.CreatedAt.Format("2006-01-02 15:04:05 -07:00"),
				c.BaseURL() + ws.RelativeURLRoot + "/" + ws.Hasher.EncodeAtomic(paste.ID),
			})
		}

		if page > 1 {
			newQuery := &fasthttp.Args{}
			c.Request().URI().QueryArgs().CopyTo(newQuery)
			newQuery.Set("page", fmt.Sprint(page-1))
			prevLink := c.BaseURL() + ws.RelativeURLRoot + "/browse?" + newQuery.String()

			t.AppendFooter(table.Row{"Prev page", prevLink})
		}
		if page < paginatedPasteList.Pagination.TotalPages {
			newQuery := &fasthttp.Args{}
			c.Request().URI().QueryArgs().CopyTo(newQuery)
			newQuery.Set("page", fmt.Sprint(page+1))
			nextLink := c.BaseURL() + ws.RelativeURLRoot + "/browse?" + newQuery.String()

			t.AppendFooter(table.Row{"Next page", nextLink})
		}

		return c.SendString(t.Render() + "\n")
	}

	return c.Render("browse/page", fiber.Map{
		"RelativeURLRoot": ws.RelativeURLRoot,
		"AbsoluteURLRoot": ws.AbsoluteURLRoot,
		"ActiveMenu":      2,
		"Action":          "browse",
		"Languages":       monaco.Languages,
		"Version":         version.Version,
		"Pastes":          paginatedPasteList.Pastes,
		"Pagination":      paginatedPasteList.Pagination,
		"Query":           c.Request().URI().QueryArgs(),
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

	paste, err := ws.DB.GetPasteByID(id)
	if err != nil {
		return err
	}

	c.Response().Header.Set("X-Filename", paste.Filename)
	c.Response().Header.Set("X-Filetype", paste.Filetype)
	return c.SendString(paste.Content)
}

// Show ...
func (ws *WebServer) Show(c *fiber.Ctx) error {
	if strings.Index(string(c.Context().UserAgent()), "curl/") == 0 {
		return ws.ShowRaw(c)
	}

	idString := c.Params("id")
	id, err := ws.Hasher.Decode(idString)
	if err != nil {
		// Mask error - we don't want to expose hasher.Decode errors
		logger.Error("%s", err)
		return httplib.NotFoundHandler(c)
	}

	paste, err := ws.DB.GetPasteByID(id)
	if err != nil {
		return err
	}

	action := "show"
	secret := c.Query("secret")
	if secret != "" && secret == paste.Secret {
		action = "edit"
	}

	return c.Render("paste/page", fiber.Map{
		"RelativeURLRoot": ws.RelativeURLRoot,
		"AbsoluteURLRoot": ws.AbsoluteURLRoot,
		"ActiveMenu":      2,
		"Action":          action,
		"Languages":       monaco.Languages,
		"Version":         version.Version,
		"Paste":           paste,
	}, "layout/layout")
}

// Update ...
func (ws *WebServer) Update(c *fiber.Ctx) error {
	idString := c.Params("id")
	id, err := ws.Hasher.Decode(idString)
	if err != nil {
		// Mask error - we don't want to expose hasher.Decode errors
		logger.Error("%s", err)
		return httplib.NotFoundHandler(c)
	}

	paste, err := ws.DB.GetPasteByID(id)
	if err != nil {
		return err
	}

	secret := c.FormValue("secret")
	if secret != paste.Secret {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	contentType := string(c.Request().Header.Peek(fiber.HeaderContentType))
	if strings.Contains(contentType, "json") {
		err := json.Unmarshal(c.Request().Body(), &paste)
		if err != nil {
			// Mask error - we don't want to expose json.Unmarshal errors
			logger.Error("%s", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}
	} else {
		f, _ := c.FormFile("f")
		if f != nil {
			paste.Filename = f.Filename

			mf, err := f.Open()
			if err != nil {
				return err
			}
			defer mf.Close()

			content, err := io.ReadAll(mf)
			if err != nil {
				return err
			}

			paste.Content = string(content)
		} else {
			paste.Content = c.FormValue("content")
		}

		if v := c.FormValue("filename"); v != "" {
			paste.Filename = v
		}

		paste.Filetype = c.FormValue("filetype")
		paste.IndentStyle = c.FormValue("indent-style")
		paste.IndentSize, _ = strconv.Atoi(c.FormValue("indent-size"))

		if c.FormValue("private") == "1" {
			paste.Private = 1
		}
	}

	err = paste.Validate()
	if err != nil {
		// Send proper error status
		c.Status(fiber.StatusUnprocessableEntity)
		return c.SendString(fmt.Sprintf("%#s", err))
	}

	paste, err = ws.DB.UpdatePaste(paste)
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

	if strings.Index(string(c.Context().UserAgent()), "curl/") == 0 {
		return c.SendString(fmt.Sprintf("%s%s/%s", c.BaseURL(), ws.RelativeURLRoot, hash))
	}

	return c.Redirect(fmt.Sprintf("%s/%s", ws.RelativeURLRoot, hash), fiber.StatusFound)
}

// Create ...
func (ws *WebServer) Create(c *fiber.Ctx) error {
	paste := model.NewPaste()

	contentType := string(c.Request().Header.Peek(fiber.HeaderContentType))
	if strings.Contains(contentType, "json") {
		err := json.Unmarshal(c.Request().Body(), &paste)
		if err != nil {
			// Mask error - we don't want to expose json.Unmarshal errors
			logger.Error("%s", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}
	} else {
		f, _ := c.FormFile("f")
		if f != nil {
			paste.Filename = f.Filename

			mf, err := f.Open()
			if err != nil {
				return err
			}
			defer mf.Close()

			content, err := io.ReadAll(mf)
			if err != nil {
				return err
			}

			paste.Content = string(content)
		} else {
			paste.Content = c.FormValue("content")
		}

		if v := c.FormValue("filename"); v != "" {
			paste.Filename = v
		}

		paste.Filetype = c.FormValue("filetype")
		paste.IndentStyle = c.FormValue("indent-style")
		paste.IndentSize, _ = strconv.Atoi(c.FormValue("indent-size"))

		if c.FormValue("private") == "1" {
			paste.Private = 1
		}
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

	if strings.Index(string(c.Context().UserAgent()), "curl/") == 0 {
		return c.SendString(c.BaseURL() + ws.RelativeURLRoot + fmt.Sprintf("/%s\n", hash))
	}

	return c.Redirect(fmt.Sprintf("%s/%s", ws.RelativeURLRoot, hash), fiber.StatusFound)
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
