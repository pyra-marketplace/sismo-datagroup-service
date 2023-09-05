package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"sync"
)

var l = IPRateLimiter{
	ips: make(map[string]int),
}

type IPRateLimiter struct {
	ips map[string]int
	mu  sync.Mutex
}

func NewRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("call into IPRateLimiter")
		// Specify the path(s) for which you want to apply rate limiting
		if c.Request.URL.Path == "/api/v1/record" {
			ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)

			l.mu.Lock()
			defer l.mu.Unlock()

			l.ips[ip]++

			if l.ips[ip] > 10 {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
				return
			}
		}

		c.Next()
	}
}
