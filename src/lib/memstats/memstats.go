package memstats

import (
	"context"
	"runtime"
	"time"

	"github.com/cihub/seelog"
)

type MemstatsCtx interface {
	Logger(string) seelog.LoggerInterface
	ConfigPeriodMemstats() time.Duration
}

type Memstats struct {
	stop chan struct{}
	done chan struct{}
	ctx  MemstatsCtx
}

func (it *Memstats) String() string              { return "[#memstats]" }
func (it *Memstats) log() seelog.LoggerInterface { return it.ctx.Logger("memstats") }
func (it *Memstats) Period() time.Duration       { return it.ctx.ConfigPeriodMemstats() }
func (it *Memstats) Stop()                       { it.stop <- struct{}{} }
func (it *Memstats) Done()                       { <-it.done }
func (it *Memstats) DoneWithContext(ctx context.Context) {
	select {
	case <-it.done:
	case <-ctx.Done():
	}
}

func (it *Memstats) GoStart() *Memstats {
	go it.Start()
	return it
}

func (it *Memstats) Start() {
	ctx := context.Background()
	it.start(ctx)
}

func (it *Memstats) StartWithCtx(ctx context.Context) {
	it.start(ctx)
}

func (it *Memstats) start(ctx context.Context) {
	if it.ctx.ConfigPeriodMemstats() == 0 {
		return
	}

	if ctx == nil { // parent-context
		ctx = context.TODO()
	}

	it.log().Infof("%s started", it)
	defer func() {
		it.log().Infof("%s finished", it)
		it.log().Flush()

		select {
		case it.done <- struct{}{}:
		case <-time.After(1 * time.Second):
		}
	}()

	ticker := time.NewTicker(it.ctx.ConfigPeriodMemstats())
	defer ticker.Stop()

	for {
		select {
		case <-it.stop:
			it.log().Infof("%s stopped", it)
			return
		case <-ctx.Done():
			it.log().Infof("%s stopped", it)
			return
		case <-ticker.C:
			it.Perform()
		}
	}
}

func (it *Memstats) Perform() {
	log := it.log()
	defer log.Flush()

	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	log.Infof(
		"(:gorutines %d :num-gc %d :alloc %d :mallocs %d :frees %d :heap-alloc %d :stack-inuse %d)",
		runtime.NumGoroutine(),
		memStats.NumGC,
		memStats.Alloc,
		memStats.Mallocs,
		memStats.Frees,
		memStats.HeapAlloc,
		memStats.StackInuse,
	)
}

func NewMemstats(ctx MemstatsCtx) *Memstats {
	it := new(Memstats)
	it.ctx = ctx
	it.stop = make(chan struct{})
	it.done = make(chan struct{})
	return it
}
