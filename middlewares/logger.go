package middlewares

import (
	"github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func NewLogger(logger *logrus.Logger, name string) negroni.Handler {
	return negronilogrus.NewMiddlewareFromLogger(logger, name)
}
