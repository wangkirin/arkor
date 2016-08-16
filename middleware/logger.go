package middleware

import (
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

var Log = logrus.New()

func Initlog(loglevel string) {
	switch loglevel {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func logger(runmode string) macaron.Handler {
	if strings.EqualFold(runmode, "dev") {
		return func(ctx *macaron.Context) {
			logrus.WithFields(logrus.Fields{
				"[Method]":        ctx.Req.Method,
				"[RequestSource]": ctx.RemoteAddr(),
				"[Time]":          time.Now().Format("2006-01-02 15:04:05"),
			}).Info(ctx.Req.RequestURI)
		}
	}
	return nil
}

// DecorateRuntimeContext appends line, file and function context to the logger
func DecorateRuntimeContext(logger *logrus.Entry) *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		return logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	} else {
		return logger
	}
}
