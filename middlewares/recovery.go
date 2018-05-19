package middlewares

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

func NewRecovery(logger logrus.FieldLogger) negroni.Handler {
	middleware := negroni.NewRecovery()
	middleware.PrintStack = false
	middleware.Logger = logger
	middleware.Formatter = &textPanicFormatter{}
	return middleware
}
