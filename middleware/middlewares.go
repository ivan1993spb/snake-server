package middleware

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type textPanicFormatter struct{}

func (t *textPanicFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *negroni.PanicInformation) {
	if rw.Header().Get("Content-Type") == "" {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	fmt.Fprintf(rw, http.StatusText(http.StatusInternalServerError))
}

func NewRecovery(logger *logrus.Logger) negroni.Handler {
	middleware := negroni.NewRecovery()
	middleware.PrintStack = false
	middleware.Logger = logger
	middleware.Formatter = &textPanicFormatter{}
	return middleware
}

const httpLoggerFormat = "request processed: {{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Method}} {{.Path}}"

func NewLogger(logger *logrus.Logger) negroni.Handler {
	middleware := negroni.NewLogger()
	middleware.SetFormat(httpLoggerFormat)
	middleware.ALogger = logger
	return middleware
}
