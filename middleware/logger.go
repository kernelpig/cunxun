package middleware

import (
	"context"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/Meiqia/melog"
	"github.com/meiqia/chi/middleware"

	"wangqingang/cunxun/account"
	"wangqingang/cunxun/utils/render"
)

const (
	AccessLogKey = "access_log"
)

func Logger(next http.Handler) http.Handler {
	return loggerMiddleware(next)
}

func loggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		dumpRequest, err := httputil.DumpRequest(r, true)
		if err != nil {
			render.InternelError(w, r, err)
			return
		}

		start := time.Now()
		al := &account.AccessLog{
			Method:   r.Method,
			Path:     r.RequestURI,
			ClientIp: r.RemoteAddr,
			Request:  string(dumpRequest),
			Response: "",
		}

		defer func() {
			al.Cost = time.Now().Sub(start).String()
			al.StatusCode = int64(ww.Status())

			melog.Info(al)
		}()

		ctx := context.WithValue(r.Context(), AccessLogKey, al)
		next.ServeHTTP(ww, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func GetAccessLog(r *http.Request) *account.AccessLog {
	return r.Context().Value(AccessLogKey).(*account.AccessLog)
}
