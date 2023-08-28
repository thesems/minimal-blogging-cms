package api

import "net/http"

func IsAuthorized(req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil || cookie.Value != "topsecret" {
		return false
	}
	return true
}
