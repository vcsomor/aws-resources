package threads

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testTask struct {
	sleep  time.Duration
	result string
}

var _ Task = (*testTask)(nil)

func (t *testTask) Execute() any {
	time.Sleep(t.sleep)
	return t.result
}

func TestThreadPool(t *testing.T) {
	const (
		threads = 5
		tasks   = 10
	)

	start := time.Now()
	var taskFutures []TaskFuture

	mgr := NewThreadpool(threads)
	for i := 0; i < tasks; i++ {
		tf, err := mgr.SubmitTask(&testTask{
			sleep:  1 * time.Second,
			result: fmt.Sprintf("result: %d", i),
		})

		assert.Nil(t, err, fmt.Sprintf("the task should not go into error %d", i))
		taskFutures = append(taskFutures, tf)
	}

	for _, tf := range taskFutures {
		resultRaw := tf.Get()
		result, ok := resultRaw.(string)
		assert.Truef(t, ok, "result is not a string")
		assert.NotNilf(t, result, "result is nil")
	}

	finished := time.Now()
	deadline := start.Add(2100 * time.Millisecond)
	assert.True(t, finished.Before(deadline), "threadpool was too slow")
}
