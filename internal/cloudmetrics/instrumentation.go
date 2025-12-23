package cloudmetrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/smallbiznis/valora/internal/config"
	"go.uber.org/fx"
)

func RegisterInstrumentation(lc fx.Lifecycle, cfg config.Config) {
	http.Handle("/metrics", promhttp.Handler())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func ()  {
				http.ListenAndServe(":2112", nil)
			}()
			return nil
		},
	})
}