package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mehgokalp/insider-project/pkg/meta"
	"time"
)

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			line := make(map[string]interface{})

			line["app_name"] = meta.AppName
			line["level"] = "debug"
			line["status_code"] = params.StatusCode
			line["path"] = params.Path
			line["method"] = params.Method
			line["remote_addr"] = params.ClientIP
			line["response_time"] = params.Latency.String()
			line["time"] = params.TimeStamp.Format(time.RFC3339)

			s, _ := json.Marshal(line)
			return string(s) + "\n"
		},
	)
}
