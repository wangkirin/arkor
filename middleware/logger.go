package middleware

import (
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

var Log = logrus.New()

func InitLog() {
	// save log content to file, Not enable yet

	// // write logs to local file
	// f, err := os.OpenFile("log/arkor.log", os.O_WRONLY|os.O_CREATE, 0755)
	// if err != nil {
	// 	fmt.Errorf("Init logger middleware FAIL: %s", err.Error())
	// }
	// logrus.SetOutput(f)

	// // set Output format
	// logrus.SetFormatter(&logrus.JSONFormatter{})
}

func logger(runmode string) macaron.Handler {
	if strings.EqualFold(runmode, "dev") {
		return func(ctx *macaron.Context) {
			Log.WithFields(logrus.Fields{
				"[Method]":        ctx.Req.Method,
				"[RequestHeader]": ctx.Req.Header,
				"[Path]":          ctx.Req.RequestURI,
			}).Info("Request Received")
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
