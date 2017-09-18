package builder

import (
	"sync"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
)

/*
func prepareState(maxWorkers int) *State {
	processor := func(logger logrus.FieldLogger, task Task) {
		time.Sleep(100 * time.Millisecond)
		task.callback(task.id, task.id)
	}

	state := NewState(processor, logrus.WithField("max_workers", maxWorkers), maxWorkers)
	return state
}

// Please, run tests with race detector: go test -v -race
func TestAddTask(t *testing.T) {
	maxWorkers := []int{1, 2, 5, 10, 20}

	for _, maxW := range maxWorkers {
		mux := &sync.RWMutex{}
		completed := make(map[string]string)
		state := prepareState(maxW)

		// callback "counts" all tasks which were processed
		callback := func(state string, description string) {
			mux.Lock()
			completed[state] = description
			mux.Unlock()
		}

		task := NewTask(callback,"1", "test", "test", "user_1", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "2", "test", "test", "user_1", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "3", "test", "test", "user_2", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "4", "test", "test", "user_3", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "5", "test", "test", "user_4", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "6", "test", "test", "user_5", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "7", "test", "test", "user_5", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "8", "test", "test", "user_5", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "9", "test", "test", "user_5", "test", "test")
		state.AddTask(&task)
		task = NewTask(callback, "10", "test", "test", "user_5", "test", "test")
		state.AddTask(&task)

		taskIDs := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
		for _, taskID := range taskIDs {
			_, ok := completed[taskID]

			if !ok {
				t.Errorf("Task %s was not completed!", taskID)
			}
		}
	}
}
*/
