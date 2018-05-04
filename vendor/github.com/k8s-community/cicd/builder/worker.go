package builder

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/cicd/builder/task"
)

// Processor is a function to process tasks
type Processor func(task task.CICD)

// worker represents a worker to process tasks
type worker struct {
	id int // unique id of the worker

	pool chan task.CICD // pool of tasks to be processed, shared between the dispatcher and the workers

	mutex      *sync.RWMutex        // mutex of inProgress queue, shared between the dispatcher and the workers
	inProgress map[string]task.CICD // queue of tasks which are currently in progress, shared between the dispatcher and the workers

	wgWorkersStopped *sync.WaitGroup // wait group to mark that workers are stoped

	cancel chan struct{} // channel to mark worker what it has to stop its work

	log       logrus.FieldLogger
	processor Processor
}

func newWorker(
	processor Processor,
	log logrus.FieldLogger,
	id int,
	pool chan task.CICD,
	mutex *sync.RWMutex,
	inProgress map[string]task.CICD,
	wgWorkersStopped *sync.WaitGroup,
) *worker {
	return &worker{
		id:               id,
		pool:             pool,
		mutex:            mutex,
		inProgress:       inProgress,
		wgWorkersStopped: wgWorkersStopped,
		processor:        processor,
		log:              log,
		cancel:           make(chan struct{}),
	}
}

// run runs the worker.
// It will work until cancellation signal is sent.
// worker reads tasks from the pool channel.
func (w *worker) run() {
	go func() {
		for {
			select {
			case t := <-w.pool:
				logger := w.log.WithField("task_id", t.ID)
				logger.Infof("worker #%d is processing task %s...", w.id, t.ID)

				w.processor(t)

				w.mutex.Lock()
				delete(w.inProgress, t.Repo)
				w.mutex.Unlock()

				logger.Infof("worker #%d processed task %s.", w.id, t.ID)

			case <-w.cancel:
				w.wgWorkersStopped.Done()
				w.log.Infof("worker #%d stopped.", w.id)
				return
			}
		}
	}()
}

// stop sends a command to stop the worker.
func (w *worker) stop() {
	w.log.Infof("worker #%d received stopping command...", w.id)
	w.cancel <- struct{}{}
}
