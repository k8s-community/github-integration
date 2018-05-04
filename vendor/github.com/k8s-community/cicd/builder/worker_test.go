package builder

import (
	"testing"

	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/cicd/builder/task"
)

func prepareWorker(pool chan task.CICD, mx *sync.RWMutex, inProgress map[string]task.CICD, wg *sync.WaitGroup) *worker {
	processor := func(taskItem task.CICD) {
		taskItem.Callback(taskItem.ID, taskItem.ID, taskItem.ID)
	}

	id := 1

	worker := newWorker(
		processor,
		logrus.WithField("id", id),
		id,
		pool,
		mx,
		inProgress,
		wg,
	)

	return worker
}

func TestRunWorker(t *testing.T) {

	pool := make(chan task.CICD)
	mx := &sync.RWMutex{}
	inProgress := make(map[string]task.CICD)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	wrk := prepareWorker(pool, mx, inProgress, wg)

	wrk.run()

	mxCompleted := &sync.RWMutex{}
	completedTasks := []string{}
	callback := func(taskID string, state string, description string) {
		mxCompleted.Lock()
		completedTasks = append(completedTasks, taskID)
		mxCompleted.Unlock()
	}

	// Add task to pool
	mx.Lock()
	taskItem := task.NewCICD(
		callback, "test", task.TypeBuild,
		"github.com", "k8s-community/cicd", "master",
		"1234", "test-namespace",
	)
	inProgress[taskItem.Repo] = *taskItem
	pool <- *taskItem
	mx.Unlock()

	wrk.stop()
	wg.Wait()

	mxCompleted.Lock()
	if len(completedTasks) == 0 {
		t.Fail()
	}

	if completedTasks[0] != "test" {
		t.Fail()
	}
	mxCompleted.Unlock()
}
