package web

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QueryString retrieves a string param from the gin request querystring
func QueryString(c *gin.Context, name string) string {
	return c.Query(name)
}

//QueryInt retrieves an integer param from the gin request querystring
func QueryInt(c *gin.Context, name string) int {
	v := c.Query(name)
	i, _ := strconv.Atoi(v)
	return i
}

//QueryDefaultInt retrieves an integer param from the gin request querystring
//defaults to def argument if not found
func QueryDefaultInteger(c *gin.Context, name string, def int) (int, error) {
	v := c.Query(name)
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

//QueryBool retrieves a boolean param from the gin request querystring
func QueryBool(c *gin.Context, name string) bool {
	return c.Query(name) == "true"
}
