package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ClientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var clients = make(map[string]*ClientLimiter)

func getLimiter(ip string) *rate.Limiter {
	if client, exists := clients[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute/5), 5) // 1 req/sec, burst 5
	clients[ip] = &ClientLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}
		c.Next()
	}
}

func StartRateLimiterCleanup() {
	go func() {
		for {
			time.Sleep(time.Minute)

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()
}
