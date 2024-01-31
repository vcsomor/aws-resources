package threads

import (
	"github.com/panjf2000/ants/v2"
	"sync/atomic"
	"time"
)

type Task interface {
	Execute() any
}

type Future interface {
	Get() any
	GetWait() any
	IsDone() bool
}

type taskFuture struct {
	getResult func() any
	isDone    func() bool
}

var _ Future = (*taskFuture)(nil)

func newFuture(getResult func() any, isDone func() bool) Future {
	return &taskFuture{
		getResult: getResult,
		isDone:    isDone,
	}
}

func (w *taskFuture) Get() any {
	return w.getResult()
}

func (w *taskFuture) GetWait() any {
	for !w.isDone() {
		time.Sleep(1 * time.Millisecond)
	}
	return w.getResult()
}

func (w *taskFuture) IsDone() bool {
	return w.isDone()
}

type Threadpool interface {
	SubmitTask(t Task) (Future, error)
	Shutdown()
}

type antsThreadPool struct {
	threadPool *ants.Pool
}

var _ Threadpool = (*antsThreadPool)(nil)

func NewThreadpool(threads int) (Threadpool, error) {
	p, err := ants.NewPool(threads, ants.WithPreAlloc(true))
	if err != nil {
		return nil, err
	}

	return &antsThreadPool{
		threadPool: p,
	}, nil
}

func (p *antsThreadPool) SubmitTask(t Task) (Future, error) {
	value := atomic.Pointer[any]{}
	done := atomic.Bool{}

	exec := func() {
		r := t.Execute()
		value.Store(&r)
		done.Store(true)
	}

	getResult := func() any {
		if l := value.Load(); l != nil {
			return *value.Load()
		}
		return nil
	}

	isDone := func() bool {
		return done.Load()
	}

	err := p.threadPool.Submit(exec)
	if err != nil {
		return nil, err
	}

	return newFuture(getResult, isDone), nil
}

func (p *antsThreadPool) Shutdown() {
	_ = p.threadPool.ReleaseTimeout(5 * time.Second)
}
