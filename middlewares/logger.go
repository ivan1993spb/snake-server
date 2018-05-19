package middlewares

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const httpLoggerFormat = "request processed: {{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Method}} {{.Path}}"

func NewLogger(logger logrus.FieldLogger) negroni.Handler {
	middleware := negroni.NewLogger()
	middleware.SetFormat(httpLoggerFormat)
	middleware.ALogger = logger
	return middleware
}
