package tracker

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			flags := map[string]string{
				"endpoint": context.Request.URL.RequestURI(),
			}
			if rval := recover(); rval != nil {
				debug.PrintStack()
				rvalStr := fmt.Sprint(rval)
				packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))
				raven.Capture(packet, flags)
				context.Writer.WriteHeader(http.StatusInternalServerError)
			}
		}()
		context.Next()
	}
}
