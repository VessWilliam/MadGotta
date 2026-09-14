package httpx

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, status int, cmp templ.Component) {
	c.Status(status)
	if err := cmp.Render(c.Request.Context(), c.Writer); err != nil {
		c.Error(err)
	}
}
