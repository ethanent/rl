package rl

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

func removePort(i string) string {
	portIdx := strings.LastIndex(i, ":")

	i = string([]byte(i)[:portIdx])

	return i
}

// LimitWrap wraps handler with a ratelimiter which counts requests
// and whose counter resets every d. Requests are rejected with a 429
// status code after counter for a single IP exceeds d.
func LimitWrap(d time.Duration, max int, handler http.HandlerFunc) http.HandlerFunc {
	counts := map[string]int{}
	countsMux := &sync.RWMutex{}

	t := time.NewTicker(d)

	go func() {
		for {
			<-t.C

			countsMux.Lock()
			counts = map[string]int{}
			countsMux.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		countsMux.Lock()
		useIP := removePort(r.RemoteAddr)
		counts[useIP]++

		curCount := counts[useIP]
		countsMux.Unlock()

		if curCount > max {
			http.Error(w, "Too many requests", 429)
			return
		}

		handler(w, r)
	}
}
