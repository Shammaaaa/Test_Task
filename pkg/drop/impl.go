package drop

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/shamil/Test_task/pkg/log"
)

type Impl struct {
	*Droppable
	context    context.Context
	cancelRoot context.CancelFunc
}

func NewContext(ctx context.Context) *Impl {
	i := &Impl{}
	i.context, i.cancelRoot = context.WithCancel(ctx)
	i.Droppable = &Droppable{}
	return i
}

func (s *Impl) Context() context.Context {
	return s.context
}

func (s *Impl) Shutdown(onError func(error)) {
	s.cancelRoot()
	log.Info("root context is canceled")
	s.EachDroppers(func(dropper Drop) {
		if impl, ok := dropper.(Debug); ok {
			log.Info(impl.DropMsg())
		}
		if err := dropper.Drop(); err != nil {
			onError(err)
		}
	})
}

func (s *Impl) Stacktrace() {
	log.Info("waiting for the server completion report to be generated")

	<-time.After(time.Second * 2)

	memory := runtime.MemStats{}
	runtime.ReadMemStats(&memory)

	colored := func(category, context string) string {
		return fmt.Sprintf("%s: %s", log.Colored(category, log.Cyan), log.Colored(context, log.Green))
	}

	fmt.Printf(
		"%s \n \t - %s \n \t - %s \n",
		log.Colored("REPORT", log.Green),
		colored("number of remaining goroutines:", strconv.Itoa(runtime.NumGoroutine())),
		colored("number of operations of the garbage collector:", strconv.Itoa(int(memory.NumGC))),
	)
}
