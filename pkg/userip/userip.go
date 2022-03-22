package userip

import "net/http"

// Get returns user IP address from given http request
func Get(r *http.Request) string {
	ips := r.Header.Values("X-Forwarded-For")
	if len(ips) == 0 {
		return r.RemoteAddr
	}

	return ips[0]
}
