package logs

import (
	"fmt"
	"io"
	"main/pkg/constants"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var ApplicationLog *logrus.Logger

func SetUpApplicationLogs() {
	logger := logrus.New()

	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logFileName := constants.ApplicationConfig.Application.LogPath + "cinema.log"
	fmt.Println("Log file ::", logFileName)

	logFile := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    100,
		MaxAge:     28,
		MaxBackups: 30000,
		LocalTime:  true,
		Compress:   true,
	}

	logger.SetOutput(io.MultiWriter(logFile))
	ApplicationLog = logger

	defer logFile.Close()

}

func InfoLog(format string, a ...any) {
	stringMessage := fmt.Sprintf(format, a...)
	ApplicationLog.WithFields(logrus.Fields{}).Info(stringMessage)
}

func ErrorLog(format string, a ...any) {
	stringMessage := fmt.Sprintf(format, a...)
	ApplicationLog.WithFields(logrus.Fields{}).Error(stringMessage)
}
