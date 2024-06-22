package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/dlbarduzzi/bookshop/internal/jsoner"
	"github.com/dlbarduzzi/bookshop/internal/logging"
)

type Limiter struct {
	// RPS is the maximum number of requests per second.
	RPS float64

	// Burst is the maximum number of burst requests.
	Burst int

	// Enable rate limit.
	Enabled bool
}

func (l Limiter) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Minute*3 {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l.Enabled {
			ctx := r.Context()
			log := logging.LoggerFromContext(ctx)

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Error("rate limit split remote address failed", "err", err)
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(l.RPS), l.Burst)}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				log.Error("rate limit exceeded", "addr", ip)
				data := jsoner.Envelope{
					"ok":         false,
					"error":      "rate limit exceeded",
					"error_code": "rate-limit-error",
				}
				if err := jsoner.Marshal(w, data, http.StatusTooManyRequests, nil); err != nil {
					w.Header().Set("Connection", "close")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}

			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}
