package jwt

import (
	"net/http"
	"strings"
)

// // extractBearerToken extracts auth token from Authorization header.
func ExtractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {

		return ""
	}

	return splitToken[1]
}
