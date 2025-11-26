package runtimeutil

import (
	"golang.org/x/sync/errgroup"
	"sync/atomic"
	"time"
)

func Go(fn func()) {
	defer func() {
		accessRecoverMu.RLock()
		defer accessRecoverMu.RUnlock()
		Recover()
	}()
	trace := GetOrCreateTrace()
	go func() {
		StoreTrace(trace)
		defer AutoRemoveTrace()
		fn()
	}()
}

func NewErrGrp() *ErrGrp {
	return &ErrGrp{
		errGrp: new(errgroup.Group),
	}
}

type ErrGrp struct {
	errGrp *errgroup.Group
}

func (e *ErrGrp) Go(fn func() error) {
	trace := GetOrCreateTrace()
	e.errGrp.Go(func() error {
		StoreTrace(trace)
		defer AutoRemoveTrace()
		return fn()
	})
}

func (e *ErrGrp) Wait() error {
	return e.errGrp.Wait()
}

type LimitedConcurErrGrp struct {
	size             uint32
	hdlCh            chan func() error
	errGrp           *ErrGrp
	quitGoWhenErr    bool
	inQueueGoCounter atomic.Int32
}

func NewLimitedConcurErrGrp(concur uint32, quitGoWhenErr bool) *LimitedConcurErrGrp {
	grp := &LimitedConcurErrGrp{
		size:          concur,
		hdlCh:         make(chan func() error, concur),
		errGrp:        NewErrGrp(),
		quitGoWhenErr: quitGoWhenErr,
	}
	return grp
}

func (g *LimitedConcurErrGrp) Go(fn func() error) {
	g.inQueueGoCounter.Add(1)
	go func() {
		g.hdlCh <- fn
	}()
}

func (g *LimitedConcurErrGrp) Wait() error {
	for i := uint32(0); i < g.size; i++ {
		g.errGrp.Go(func() error {
			var lastErr error
			for {
				var hdl func() error
				select {
				case hdl = <-g.hdlCh:
					g.inQueueGoCounter.Add(-1)
					if err := hdl(); err != nil {
						if g.quitGoWhenErr {
							// 先清空channel中还没完成的hdl在退出
							for {
								select {
								case <-g.hdlCh:
									g.inQueueGoCounter.Add(-1)
								default:
									return err
								}
							}
						}
						lastErr = err
					}
				default:
					// 确保例程写入hdl的线程已经全部写完
					if g.inQueueGoCounter.Load() > 0 {
						time.Sleep(time.Millisecond)
						continue
					}
					goto out
				}
			}
		out:
			return lastErr
		})
	}

	return g.errGrp.Wait()
}
