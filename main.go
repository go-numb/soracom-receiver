package main

import (
	"io"
	"net/http"
	"text/template"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	PORT = ":8080"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Use(
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*", "http://localhost:3000"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-SERVICE-KEY"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		}))
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	e.Static("/static", "static")

	// HTML handlers
	e.GET("/", handlerToHTMLRender)

	// API handlers
	e.GET("/api/v1/receive", handlerForQuery)
	e.POST("/api/v1/receive", handlerForParam)

	e.Logger.Fatal(e.Start(PORT))
}

type Request struct {
	// Type 処理分岐用
	Type int
	Data interface{}

	Timestamp int64
}

type Response struct {
	Code      string      `json:"code"`
	Error     string      `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

// handlerToHTMLRender HTML Renderer
func handlerToHTMLRender(e echo.Context) error {
	return e.Render(http.StatusOK, "layout", nil)
}

// handlerForQuery
func handlerForQuery(e echo.Context) error {
	r := &Response{
		Code:      "success",
		Timestamp: time.Now(),
	}
	return e.JSON(http.StatusOK, r)
}

// handlerToParam 取得データ加工・通知用
func handlerForParam(e echo.Context) error {
	r := new(Request)
	if err := e.Bind(r); err != nil {
		return e.JSON(http.StatusBadRequest, nil)
	}

	switch r.Type {
	case 0:
		// dosomething()
	case 1:
		// dosomething()
	default:
	}

	res := &Response{
		Code:      "success",
		Timestamp: time.Now(),
	}
	return e.JSON(http.StatusOK, res)
}
