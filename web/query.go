package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryString(c *gin.Context, name string) string {
	return c.Query(name)
}

func QueryInt(c *gin.Context, name string) int {
	v := c.Query(name)
	i, _ := strconv.Atoi(v)
	return i
}

func QueryBool(c *gin.Context, name string) bool {
	return c.Query(name) == "true"
}
