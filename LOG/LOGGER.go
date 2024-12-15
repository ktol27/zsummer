package LOG

import (
	"github.com/sirupsen/logrus"
	"os"
)

var ErrorLogger *logrus.Logger

func InitLogger() {
	ErrorLogger = logrus.New()
	ErrorLogger.SetReportCaller(true)
	ErrorLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
	ErrorLogger.SetOutput(os.Stdout)
	ErrorLogger.SetLevel(logrus.ErrorLevel)
}
