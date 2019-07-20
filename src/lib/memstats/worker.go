package memstats

import (
	"context"
	"fmt"
	"runtime"
	"time"

	bufPool "github.com/vany-egorov/ha-eta/lib/buf-pool"
	"github.com/vany-egorov/ha-eta/lib/log"
)

type Worker struct {
	stop chan struct{}
	done chan struct{}
	cfg  Config
}

func (it *Worker) String() string { return "(:memstats)" }
func (it *Worker) Stop()          { it.stop <- struct{}{} }
func (it *Worker) Done()          { <-it.done }
func (it *Worker) DoneWithContext(ctx context.Context) {
	select {
	case <-it.done:
	case <-ctx.Done():
	}
}

func (it *Worker) Start() {
	ctx := context.Background()
	it.start(ctx)
}

func (it *Worker) StartWithCtx(ctx context.Context) {
	it.start(ctx)
}

func (it *Worker) start(ctx context.Context) {
	if ctx == nil { // parent-context
		ctx = context.TODO()
	}

	it.cfg.fnLog(log.Info,
		fmt.Sprintf("%s started", it))
	defer func() {
		it.cfg.fnLog(log.Info,
			fmt.Sprintf("%s finished", it))

		select {
		case it.done <- struct{}{}:
		default:
		}
	}()

	ticker := time.NewTicker(it.cfg.period)
	defer ticker.Stop()

	for {
		select {
		case <-it.stop:
			it.cfg.fnLog(log.Info,
				fmt.Sprintf("%s stoped", it))
			return
		case <-ctx.Done():
			it.cfg.fnLog(log.Info,
				fmt.Sprintf("%s stoped", it))
			return
		case <-ticker.C:
			it.Perform()
		}
	}
}

func (it *Worker) Perform() {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	buf := bufPool.NewBuf()
	defer buf.Release()

	buf.WriteString(it.String())
	buf.WriteByte(' ')
	buf.WriteString(fmt.Sprintf("(:gorutines %d :num-gc %d :alloc %d :mallocs %d :frees %d :heap-alloc %d :stack-inuse %d)",
		runtime.NumGoroutine(),
		memStats.NumGC,
		memStats.Alloc,
		memStats.Mallocs,
		memStats.Frees,
		memStats.HeapAlloc,
		memStats.StackInuse,
	))

	it.cfg.fnLog(log.Info, buf.String())
}

func (it *Worker) Initialize(fnArgs ...Arg) {
	it.stop = make(chan struct{})
	it.done = make(chan struct{}, 1)

	it.cfg.defaultize()

	for _, fn := range fnArgs {
		fn(&it.cfg)
	}
}
