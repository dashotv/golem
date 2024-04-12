package router

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

// QueryParamString retrieves a string param from the gin request querystring
func QueryParamString(c echo.Context, name string) string {
	return c.QueryParam(name)
}

// QueryParamInt retrieves an integer param from the gin request querystring
func QueryParamInt(c echo.Context, name string) int {
	v := c.QueryParam(name)
	i, _ := strconv.Atoi(v)
	return i
}

// QueryBool retrieves a boolean param from the gin request querystring
func QueryParamBool(c echo.Context, name string) bool {
	return c.QueryParam(name) == "true"
}

// QueryDefaultInt retrieves an integer param from the gin request querystring
// defaults to def argument if not found
func QueryDefaultInt(c echo.Context, name string, def int) (int, error) {
	v := c.QueryParam(name)
	if v == "" {
		return def, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return def, err
	}

	if n < 0 {
		return def, fmt.Errorf("less than zero")
	}

	return n, nil
}

func QueryParamIntDefault(c echo.Context, name string, def string) int {
	param := c.QueryParam(name)
	result, err := strconv.Atoi(param)
	if err == nil && result > 0 {
		return result
	}

	d, err := strconv.Atoi(def)
	if err == nil && d > 0 {
		return d
	}

	return 0
}

func QueryParamFloatDefault(c echo.Context, name string, def string) float64 {
	param := c.QueryParam(name)
	result, err := strconv.ParseFloat(param, 64)
	if err == nil && result > 0 {
		return result
	}

	if def != "" {
		d, err := strconv.ParseFloat(param, 64)
		if err == nil && d > 0 {
			return d
		}
	}

	return 0
}

func QueryParamBoolDefault(c echo.Context, name string, def string) bool {
	param := c.QueryParam(name)
	if param != "" {
		return param == "true"
	}

	if def != "" {
		return def == "true"
	}

	return false
}

func QueryParamStringDefault(c echo.Context, name string, def string) string {
	param := c.QueryParam(name)
	if param != "" {
		return param
	}
	return def
}

// stolen from gin gonic
// H is a shortcut for map[string]any
type H map[string]any

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
