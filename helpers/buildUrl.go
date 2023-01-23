package helper

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BuildURL(c *gin.Context, page int, limit int) string {
    req := c.Request
    u, _ := url.Parse(req.URL.String())
    q := u.Query()

    var protocol string
    if req.TLS != nil {
        protocol = "https"
    } else {
        protocol = "http"
    }

    host := protocol + "://" + req.Host
    q.Set("page", strconv.Itoa(page))
    q.Set("limit", strconv.Itoa(limit))
    u.RawQuery = q.Encode()
    return host + u.String()
}