package tracker

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			meta := map[string]string{
				"endpoint": context.Request.URL.RequestURI(),
			}
			if err := recover(); err != nil {
				TrackError(err, "TopLevelError", meta)
				context.Writer.WriteHeader(http.StatusInternalServerError)
			}
		}()
		context.Next()
	}
}
