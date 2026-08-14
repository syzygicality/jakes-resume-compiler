package middleware

import (
	"fmt"
	"jakes-resume-compiler/server/platform/utils"
	"net/http"
)

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = fmt.Errorf("%v", rec)
				}
				utils.HTTPError(err, w, r, "panic recovered", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
