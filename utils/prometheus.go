package utils

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8080", nil)
	}()
}
