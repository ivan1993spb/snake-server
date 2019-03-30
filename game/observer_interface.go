package game

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/world"
)

// TODO: Delete.
type ObserverInterface interface {
	Observe(stop <-chan struct{}, world world.Interface, logger logrus.FieldLogger)
}
