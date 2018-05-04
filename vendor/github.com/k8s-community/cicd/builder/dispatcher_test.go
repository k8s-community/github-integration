package builder

import (
	"sync"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/cicd/builder/task"
)

func prepareDispatcher(maxWorkers int) *Dispatcher {
	processor := func(taskItem task.CICD) {
		taskItem.Callback(taskItem.ID, taskItem.ID, taskItem.ID)
	}

	disp := NewDispatcher(processor, logrus.WithField("max_workers", maxWorkers), maxWorkers, 500*time.Millisecond)
	return disp
}

// Please, run tests with race detector: go test -v -race
func TestAddTask(t *testing.T) {
	maxWorkers := []int{1, 2, 5, 10, 20}

	for _, maxW := range maxWorkers {
		mux := &sync.RWMutex{}
		completed := make(map[string]string)
		disp := prepareDispatcher(maxW)

		// callback "counts" all tasks which were processed
		callback := func(taskID string, state string, description string) {
			mux.Lock()
			completed[taskID] = description
			mux.Unlock()
		}

		taskItem := task.NewCICD(callback, "1", "test", "test", "user_1", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "2", "test", "test", "user_1", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "3", "test", "test", "user_2", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "4", "test", "test", "user_3", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "5", "test", "test", "user_4", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "6", "test", "test", "user_5", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "7", "test", "test", "user_5", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "8", "test", "test", "user_5", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "9", "test", "test", "user_5", "test", "test", "test-namespace")
		disp.AddTask(taskItem)
		taskItem = task.NewCICD(callback, "10", "test", "test", "user_5", "test", "test", "test-namespace")
		disp.AddTask(taskItem)

		disp.Shutdown()

		taskIDs := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
		for _, taskID := range taskIDs {
			mux.Lock()
			_, ok := completed[taskID]
			mux.Unlock()

			if !ok {
				t.Errorf("Task %s was not completed!", taskID)
			}
		}
	}
}
