package middlewares

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const httpLoggerFormat = "request processed: {{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Method}} {{.Path}}"

// TODO: Use https://github.com/meatballhat/negroni-logrus ?

func NewLogger(logger logrus.FieldLogger) negroni.Handler {
	middleware := negroni.NewLogger()
	middleware.SetFormat(httpLoggerFormat)
	middleware.ALogger = logger
	return middleware
}
