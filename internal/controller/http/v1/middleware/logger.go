package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosuit/gins"
	"github.com/gosuit/lec"
	"github.com/gosuit/sl"
)

func (m *Middleware) InitLogger(c lec.Context) gin.HandlerFunc {
	log := c.Logger()

	log.Info("logger middleware enabled.")

	return func(c *gin.Context) {
		ctx := lec.New(log)

		c.Set(gins.CtxKey, ctx)

		req := c.Request

		c.Next()

		l := log.ToSlog()

		entry := l.With(
			sl.StringAttr("method", req.Method),
			sl.StringAttr("path", req.URL.Path),
			sl.StringAttr("remote_addr", req.RemoteAddr),
			sl.StringAttr("user_agent", req.UserAgent()),
		)

		t1 := time.Now()
		defer func() {
			entry.Info("request completed",
				sl.IntAttr("status", c.Writer.Status()),
				sl.StringAttr("duration", time.Since(t1).String()),
			)
		}()
	}
}
