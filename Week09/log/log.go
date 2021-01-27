package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	f, err := os.OpenFile(os.Getenv("LOG")+"services.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})
}

func Info(args ...interface{}) {
	log.Infoln(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Error(args ...interface{}) {
	log.Errorln(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
