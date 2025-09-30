package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// RequestLogMiddleware 记录每个 HTTP 请求的 Method、Path、状态码与耗时
func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		logx.Infof("HTTP %s %s status=%d latency=%s ua=%s ip=%s", r.Method, r.URL.RequestURI(), rec.status, duration, r.UserAgent(), r.RemoteAddr)
	})
}
