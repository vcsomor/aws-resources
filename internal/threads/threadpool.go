package threads

import "github.com/shettyh/threadpool"

const queueSize = 1000

// =================== TASK ====================================

type Task interface {
	Execute() any
}

type taskCallableWrapper struct {
	task Task
}

var _ threadpool.Callable = (*taskCallableWrapper)(nil)

func wrapToCallable(t Task) threadpool.Callable {
	return &taskCallableWrapper{
		task: t,
	}
}

func (w *taskCallableWrapper) Call() any {
	return w.task.Execute()
}

// =================== FUTURE ====================================

type TaskFuture interface {
	Get() any
	IsDone() bool
}

type taskFutureWrapper struct {
	future *threadpool.Future
}

var _ TaskFuture = (*taskFutureWrapper)(nil)

func wrapToTaskFuture(f *threadpool.Future) TaskFuture {
	return &taskFutureWrapper{
		future: f,
	}
}

func (w *taskFutureWrapper) Get() any {
	return w.future.Get()
}

func (w *taskFutureWrapper) IsDone() bool {
	return w.future.IsDone()
}

// =================== MANAGER ====================================

type Threadpool interface {
	SubmitTask(t Task) (TaskFuture, error)
	Shutdown()
}

type defaultThreadpool struct {
	threadPool *threadpool.ThreadPool
}

var _ Threadpool = (*defaultThreadpool)(nil)

func NewThreadpool(threads int) Threadpool {
	return &defaultThreadpool{
		threadPool: threadpool.NewThreadPool(threads, queueSize),
	}
}

func (p *defaultThreadpool) SubmitTask(t Task) (TaskFuture, error) {
	f, err := p.threadPool.ExecuteFuture(wrapToCallable(t))
	if err != nil {
		return nil, err
	}
	return wrapToTaskFuture(f), nil
}

func (p *defaultThreadpool) Shutdown() {
	p.threadPool.Close()
}
