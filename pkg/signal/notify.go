package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func Notifier(on ...func()) (wait func() error, stop func(err ...error)) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	err := make(chan error, 1)

	wait = func() error {
		<-sig
		if len(on) > 0 {
			on[0]()
		}
		select {
		case e := <-err:
			return e
		default:
			return nil
		}
	}

	stop = func(e ...error) {
		if len(e) > 0 {
			err <- e[0]
		}
		sig <- syscall.SIGINT
	}

	return wait, stop
}
