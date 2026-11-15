package logger

import (
	"fmt"
	"time"
)

func CreateOperationLogs(operationName string, infoParameters ...string) {
	logMaps := map[string]any{
		"operation": operationName,
		"time":      time.Now(),
	}

	i := 0
	for i < len(infoParameters) {
		key := fmt.Sprintf(infoParameters[i])
		if len(infoParameters) == 1 {
			break
		}
		logMaps[key] = infoParameters[i+1]

		i += 2
	}

	CreateInfoLog(logMaps, "Operation completed")
}

func CreateInMemoryCacheOperationInfoLog(cacheName string) {
	CreateOperationLogs("loadServerCache", "cacheName", cacheName)
}
