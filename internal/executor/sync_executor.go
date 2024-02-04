package executor

type SynchronousResult struct {
	Outcome any
	Error   error
}

type SynchronousExecutor interface {
	Execute(task Task) (any, error)
	ExecuteAll(tasks []Task) []SynchronousResult
}

type syncExecutionData struct {
	future  Future
	outcome any
	err     error
}

type syncExecutor struct {
	p Threadpool
}

func NewSynchronousExecutor(p Threadpool) SynchronousExecutor {
	return &syncExecutor{
		p: p,
	}
}

var _ SynchronousExecutor = (*syncExecutor)(nil)

func (e *syncExecutor) Execute(task Task) (any, error) {
	f, err := e.p.SubmitTask(task)
	if err != nil {
		return nil, err
	}
	return f.GetWait(), nil
}

func (e *syncExecutor) ExecuteAll(tasks []Task) []SynchronousResult {
	var r []*syncExecutionData

	for _, t := range tasks {
		f, err := e.p.SubmitTask(t)
		if err != nil {
			r = append(r, &syncExecutionData{err: err})
			continue
		}
		r = append(r, &syncExecutionData{future: f})
	}

	for _, d := range r {
		currentDt := d
		f := currentDt.future
		if d.err != nil || f == nil {
			continue
		}

		d.outcome = f.GetWait()
	}
	return transformSyncResult(r)
}

func transformSyncResult(r []*syncExecutionData) (results []SynchronousResult) {
	for _, dt := range r {
		results = append(results, SynchronousResult{
			Outcome: dt.outcome,
			Error:   dt.err,
		})
	}
	return results
}
