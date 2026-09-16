package utils

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// SelfPing periodically hits the app's own health endpoint to stop Render's
// free-tier web services from spinning down after 15 minutes of inactivity.
// MVP-only workaround; remove once the service is on a paid/always-on plan.
func SelfPing(appURL string, stop <-chan struct{}) {
	target := strings.TrimSuffix(appURL, "/") + "/ping"
	client := &http.Client{Timeout: 10 * time.Second} // never hang a background goroutine on a stalled connection

	ticker := time.NewTicker(14 * time.Minute)
	defer ticker.Stop()

	log.Printf("[self-ping] armed for %s every 14m (MVP-only Render free-tier workaround)", target)

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			resp, err := client.Get(target)
			if err != nil {
				log.Printf("[self-ping] failed: %v", err)
				continue
			}
			resp.Body.Close()
			log.Printf("[self-ping] ok (status %d)", resp.StatusCode)
		}
	}
}
