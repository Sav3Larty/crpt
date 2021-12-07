package internal

import (
	"context"
	"fmt"
	"forum/dbs"
	"forum/models"
	"log"
	"net"
	"net/http"
	"time"
)

type ctxKey int8

const ctxKeyUser ctxKey = iota

//Middleware for logging time for request execution and authorization.
func Middleware(limiter *models.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		limiter.Lock()
		_, ok := limiter.IPs[ip]
		if !ok {
			limiter.IPs[ip] = &models.Counter{Count: 0, LastSeen: time.Now()}
		}
		if limiter.IPs[ip].Count > 30 {
			limiter.Unlock()
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		limiter.IPs[ip].Count++
		limiter.IPs[ip].LastSeen = time.Now()
		time.AfterFunc(10*time.Second, func() {
			limiter.Lock()
			limiter.IPs[ip].Count--
			limiter.Unlock()
		})
		req := fmt.Sprintf("%s %s", r.Method, r.URL)

		user := &models.User{Username: "Guest"}
		c, err := r.Cookie("session")
		if err != http.ErrNoCookie {
			user, _ = dbs.FindUserBySession(c.Value)
		}
		limiter.Unlock()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
		log.Println(req, "completed in", time.Since(start), ip)
	})
}

func CleanupVisitors(limiter *models.Limiter) {
	for {
		time.Sleep(time.Minute)
		limiter.Lock()
		for ip := range limiter.IPs {
			if time.Since(limiter.IPs[ip].LastSeen) > 3*time.Minute {
				delete(limiter.IPs, ip)
			}
		}
		limiter.Unlock()
	}
}
