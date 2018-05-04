package builder

import (
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/cicd/builder/task"
)

// Dispatcher is a state of building process, it must live during the service is working
type Dispatcher struct {
	logger logrus.FieldLogger

	workers []*worker      // workers to process the tasks
	pool    chan task.CICD // Tasks processing by workers

	mxQueues    *sync.RWMutex        // Mutex to protect inProgress and waiting queues
	inProgress  map[string]task.CICD // Tasks which are currently in progress
	waiting     []task.CICD          // Tasks which are waiting for free worker
	reassigning map[string]task.CICD

	waitingQueueReady chan struct{} // Event marking what waiting queue is not empty

	mxShuttingDown *sync.RWMutex
	isShuttingDown bool // Flag marking what state has to stop working

	wgWorkersStopped *sync.WaitGroup // WaitGroup is waiting for workers stopping

	medianProcessorExecTime time.Duration
}

// NewDispatcher creates new Dispatcher instance with the parameters:
// - processor to process a task (use builder.Process for it)
// - logger to log important information
// - special channel to be able to process only workers goroutines in the same time
// - list of processing tasks and a mutex to deal with them
// - list of 'to do' tasks and a mutex to deal with them
// - shutdown channel to mark that service is not available for tasks anymore
func NewDispatcher(processor Processor, logger logrus.FieldLogger, maxWorkers int, waitBeforeReassign time.Duration) *Dispatcher {
	state := &Dispatcher{
		logger: logger,

		pool: make(chan task.CICD, maxWorkers),

		mxQueues:   &sync.RWMutex{},
		inProgress: make(map[string]task.CICD),

		waitingQueueReady: make(chan struct{}),

		reassigning: make(map[string]task.CICD),

		mxShuttingDown: &sync.RWMutex{},
		isShuttingDown: false,

		wgWorkersStopped: &sync.WaitGroup{},

		medianProcessorExecTime: waitBeforeReassign,
	}

	state.wgWorkersStopped.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		worker := newWorker(processor, logger, i, state.pool, state.mxQueues, state.inProgress, state.wgWorkersStopped)
		worker.run()

		state.workers = append(state.workers, worker)
	}

	go state.processWaitingQueue()

	return state
}

// Shutdown makes 'graceful shutdown' of the dispatcher and the workers:
// - Tasks, which are currently processing by the workers, will be finished
// - New tasks will be rejected
// - Scheduled tasks will be canceled
func (state *Dispatcher) Shutdown() {
	state.logger.Info("Dispatcher is preparing for shutdown...")

	state.logger.Info("Mark that service is shutting down and couldn't receive tasks anymore...")
	state.mxShuttingDown.Lock()
	state.isShuttingDown = true
	state.mxShuttingDown.Unlock()
	state.logger.Info("Marked.")

	state.logger.Info("Send 'stop' signals to workers...")
	state.logger.Info()
	for _, worker := range state.workers {
		worker.stop()
	}
	state.logger.Info("Signals were sent.")

	state.logger.Info("Wait until workers stopped...")
	state.wgWorkersStopped.Wait()
	state.logger.Info("Done.")

	state.logger.Info("Check if there are no tasks in the 'waiting' queue")
	state.mxQueues.Lock()
	for _, taskItem := range state.waiting {
		msg := state.cancelTask(&taskItem)
		state.logger.Info("Shutdown of 'waiting' queue:" + msg)
	}
	state.waiting = state.waiting[:0]
	state.mxQueues.Unlock()
	state.logger.Info("Done.")

	state.logger.Info("Check if there are no tasks in the 'inProgress' queue")
	state.mxQueues.Lock()
	for repo, taskItem := range state.inProgress {
		msg := state.cancelTask(&taskItem)
		state.logger.Info("Shutdown of 'in progress' queue: " + msg)
		delete(state.inProgress, repo)
	}
	state.mxQueues.Unlock()
	state.logger.Info("Done")

	state.logger.Info("Check if there are no tasks in the 'reAssign' queue...")
	for {
		state.mxQueues.Lock()
		length := len(state.reassigning)
		if length > 0 {
			state.mxQueues.Unlock()
			time.Sleep(state.medianProcessorExecTime)
		} else {
			state.mxQueues.Unlock()
			break
		}
	}
	state.logger.Info("Done.")

	state.logger.Info("Dispatcher was shutdown successfully.")
}

// AddTask add Task to the pool.
func (state *Dispatcher) AddTask(t *task.CICD) error {
	// check if service is shutting down
	state.mxShuttingDown.Lock()
	if state.isShuttingDown {
		// we couldn't process task because service is shutting down
		msg := state.cancelTask(t)
		state.logger.Info("AddTask: " + msg)

		state.mxShuttingDown.Unlock()
		return fmt.Errorf(msg)
	}
	state.mxShuttingDown.Unlock()

	logger := state.logger.WithField("task_id", t.ID)
	logger.Infof("Add task %s to the waiting queue...", t.ID)

	state.mxQueues.Lock()
	state.waiting = append(state.waiting, *t)
	state.mxQueues.Unlock()

	state.waitingQueueReady <- struct{}{}
	logger.Info("Task added.")

	return nil
}

// GetTasks get Tasks to the pool.
func (state *Dispatcher) GetTasks() (queue []string, progress []string, reassign []string) {
	state.mxQueues.Lock()
	defer state.mxQueues.Unlock()

	for _, item := range state.waiting {
		queue = append(queue, item.Prefix+"/"+item.Repo+": "+item.ID)
	}

	for _, item := range state.inProgress {
		progress = append(progress, item.Prefix+"/"+item.Repo+": "+item.ID)
	}

	for _, item := range state.reassigning {
		reassign = append(reassign, item.Prefix+"/"+item.Repo+": "+item.ID)
	}

	return queue, progress, reassign
}

// processWaitingQueue gets task from the "waiting" queue.
// For each repo only one task might be processing in the same time, so if the "inProgress" queue already contains
// that repo, current task will not be added and will be moved to the end of the "waiting" queue.
func (state *Dispatcher) processWaitingQueue() {
	for {
		select {
		case <-state.waitingQueueReady:
			state.logger.Info("WaitingQueue processor has caught what waiting queue is ready to be processed...")

			state.mxQueues.Lock()

			// if there are no waiting tasks because of shutting down, just continue
			if len(state.waiting) == 0 {
				state.mxQueues.Unlock()
				continue
			}

			t := state.waiting[0]

			logger := state.logger.WithField("task_id", t.ID)
			logger.Debugf("Task %s is getting from the waiting queue...", t.ID)

			inProgress, ok := state.inProgress[t.Repo]

			// Define if we can move task to the "in progress" queue
			// (we can do it only if current repo is not in the "in progress" queue yet)
			addToInProgress := false
			if !ok || (inProgress.Repo != t.Repo) {
				state.inProgress[t.Repo] = t
				addToInProgress = true
				logger.Debugf("Task %s is moving to the 'in progress' queue and is going to be processed...", t.ID)
			}

			state.waiting = state.waiting[1:]

			state.mxQueues.Unlock()

			// Send current task to the workers or move it back to the "waiting" queue
			if addToInProgress {
				state.pool <- t
			} else {
				go func() {
					// Task processing takes some time, so to not to repeat the process too many times,
					// just wait a little before re-add task to the waiting queue
					state.mxQueues.Lock()
					state.reassigning[t.ID] = t
					state.mxQueues.Unlock()

					time.Sleep(state.medianProcessorExecTime)

					state.AddTask(&t)

					state.mxQueues.Lock()
					delete(state.reassigning, t.ID)
					state.mxQueues.Unlock()
				}()
				logger.Debugf("Task %s moved back to the 'waiting' queue.", t.ID)
			}
		}
	}
}

func (state *Dispatcher) queueLen() int {
	state.mxQueues.Lock()
	defer state.mxQueues.Unlock()

	return len(state.waiting)
}

func (state *Dispatcher) cancelTask(t *task.CICD) string {
	t.Callback(t.ID, task.StateError, "CI/CD service is shutting down, please, try again later.")
	return fmt.Sprintf("Task %s wasn't processed because the service is shutting down", t.ID)
}
