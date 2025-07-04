package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Simple in-memory rate limiter (per IP)
var visitors = make(map[string]*visitor)
var mu sync.Mutex

type visitor struct {
	lastSeen time.Time
	limiter  *time.Ticker
	count    int
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func getVisitor(ip string) *visitor {
	mu.Lock()
	defer mu.Unlock()
	v, exists := visitors[ip]
	if !exists {
		v = &visitor{lastSeen: time.Now(), limiter: time.NewTicker(time.Minute), count: 0}
		visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v
}

// RateLimitMiddleware limita a 60 requisições por minuto por IP
func RateLimitMiddleware() gin.HandlerFunc {
	go cleanupVisitors()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		v := getVisitor(ip)
		if v.count >= 60 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Limite de requisições atingido. Tente novamente em instantes."})
			return
		}
		v.count++
		c.Next()
	}
}
