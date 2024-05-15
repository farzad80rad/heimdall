package proxy

import (
	"github.com/gin-gonic/gin"
)

type Proxy interface {
	Proxy(ctx *gin.Context) error
	Ping(url string) bool
}
