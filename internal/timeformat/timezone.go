package timeformat

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// this middlewares function is to save the timezone to the other handlers
type timezoneContextKey struct{}

func Timezone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		loc := time.UTC // this is a timezone, not a time. type *time.Location
		// compare to time.Now().UTC(), this is of type time.Time
		cookie, err := r.Cookie("tz")
		if err == nil {
			userLocation, err := time.LoadLocation(cookie.Value)
			if err == nil {
				loc = userLocation
			}
		}
		ctx = context.WithValue(ctx, timezoneContextKey{}, loc)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetTimezone(ctx context.Context) *time.Location {
	timezone, ok := ctx.Value(timezoneContextKey{}).(*time.Location)
	if !ok {
		slog.Error("timezone not found in the context")
		return time.UTC
	}
	return timezone
}
