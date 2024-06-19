package httputils

import (
	"net/http"
	"strings"
)

func RequestIP(r *http.Request) string {
	split := strings.Split(r.RemoteAddr, ":")
	return strings.Join(split[:len(split)-1], ":")
}
