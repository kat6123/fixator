package handler

import (
	"net/http"
	"time"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

func inPeriod(startMinutes, endMinutes TimePeriod) middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			nowMinutes := now.Hour()*60 + now.Minute()

			if int(startMinutes) <= nowMinutes && nowMinutes <= int(endMinutes) {
				f(w, r)
			}

			writeJSON(w, "make request later", http.StatusServiceUnavailable)
		}
	}
}
