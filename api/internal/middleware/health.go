package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lardira/playtrack/internal/tech"
)

func HealthCheck(api huma.API, checker *tech.HealthChecker) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if !checker.Ok() {
			log.Printf("request denied on health check: on %q", ctx.URL().Path)
			huma.WriteErr(api, ctx, http.StatusServiceUnavailable,
				"service currently unavailable", fmt.Errorf("error on a system check"),
			)
			return
		}
		next(ctx)
	}
}
