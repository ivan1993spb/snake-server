package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type responsePanic struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type textPanicFormatter struct {
	logger logrus.FieldLogger
}

// Implement PanicFormatter interface
func (t *textPanicFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *negroni.PanicInformation) {
	t.writeResponseJSON(rw, http.StatusInternalServerError, &responsePanic{
		Code: http.StatusInternalServerError,
		Text: "panic occurred",
	})
}

func (t *textPanicFormatter) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.logger.WithError(err).Error("cannot send response on panic")
	}
}

func NewRecovery(logger logrus.FieldLogger) negroni.Handler {
	middleware := negroni.NewRecovery()
	middleware.PrintStack = false
	middleware.Logger = logger
	middleware.Formatter = &textPanicFormatter{
		logger: logger,
	}
	return middleware
}
