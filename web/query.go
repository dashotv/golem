package web

import (
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

//QueryBool retrieves a boolean param from the gin request querystring
func QueryBool(c *gin.Context, name string) bool {
	return c.Query(name) == "true"
}
