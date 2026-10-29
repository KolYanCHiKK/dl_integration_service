package logger

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type DefaultErrorLogsParams struct {
	Error error
	Info  string
}

func CreateErrLog(err string, info string) {
	log.WithFields(map[string]any{
		"error": err,
	}).Error(info)
}

func createFieldsLog(fields log.Fields, info string) {
	log.WithFields(fields).Error(info)
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

func CreateValidationErrLogs(integrationType string, errParameters ...string) {
	logMaps := map[string]any{
		"integrationType": integrationType,
		"time":            time.Now(),
	}

	if len(errParameters) == 0 {
		return
	}

	i := 0
	for i < len(errParameters) {
		key := fmt.Sprintf("error[%d]", i)
		logMaps[key] = errParameters[i]

		i += 1
	}

	createFieldsLog(logMaps, "Validation is failed")
}
