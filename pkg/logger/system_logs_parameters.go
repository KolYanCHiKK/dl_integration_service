package logger

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func CreateErrLog(err string, info string) {
	log.WithFields(map[string]any{
		"error": err,
	}).Error(info)
}

func CreateInfoLog(fields log.Fields, info string) {
	log.WithFields(fields).Info(info)
}

func CreateWarnLog(fields log.Fields, info string) {
	log.WithFields(fields).Warn(info)
}

func CreateDebugLog(fields log.Fields, info string) {
	log.WithFields(fields).Debug(info)
}

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
