package rl

import (
	"net/http"
	"sync"
	"time"
)

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
		counts[r.RemoteAddr]++

		curCount := counts[r.RemoteAddr]
		countsMux.Unlock()

		if curCount > max {
			http.Error(w, "Too many requests", 429)
			return
		}

		handler(w, r)
	}
}
