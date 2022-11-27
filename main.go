package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// functions used for funcmap
func A(s string) string {
	return s
}
func Add(a, b int) int {
	c := a + b
	return c
}
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

type TemplateRenderer struct {
	templates *template.Template
}

func main() {
	e := echo.New()
	//setting the funcmaps
	renderer := &TemplateRenderer{
		templates: template.Must(template.New("t").Funcs(template.FuncMap{
			"A":   A,
			"Add": Add,
		}).ParseGlob("views/*.html")),
	}
	e.Renderer = renderer
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	// Named route "foobar"
	e.GET("/something", func(c echo.Context) error {
		return c.Render(http.StatusOK, "template.html", map[string]interface{}{
			"name": "Dolly!",
		})
	}).Name = "foobar"

	e.Logger.Fatal(e.Start(":8000"))
}
