package game

import (
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

func countdown(cxt context.Context, delay time.Duration, start int,
) <-chan int {

	if err := cxt.Err(); err != nil {
		return nil /*, fmt.Errorf("cannot start counter:", err)*/
	}
	var output = make(chan int)

	glog.Infoln("starting countdown")

	go func() {
		var ticker = time.NewTicker(delay)

		defer func() {
			close(output)
			ticker.Stop()
		}()

		defer glog.Infoln("finishing countdown")

		for ; start >= 0; start-- {
			select {
			case <-cxt.Done():
				return
			case <-ticker.C:
				output <- start
			}
		}
	}()

	return output /*, nil*/
}
