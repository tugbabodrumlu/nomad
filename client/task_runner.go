package client

import (
	"log"
	"sync"

	"github.com/hashicorp/nomad/nomad/structs"
)

// TaskRunner is used to wrap a task within an allocation and provide the execution context.
type TaskRunner struct {
	ctx    *AllocRunner
	logger *log.Logger

	task *structs.Task

	updateCh chan *structs.Task

	destroy     bool
	destroyCh   chan struct{}
	destroyLock sync.Mutex
}

// NewTaskRunner is used to create a new task context
func NewTaskRunner(ctx *AllocRunner, task *structs.Task) *TaskRunner {
	tc := &TaskRunner{
		ctx:       ctx,
		logger:    ctx.logger,
		task:      task,
		updateCh:  make(chan *structs.Task, 8),
		destroyCh: make(chan struct{}),
	}
	return tc
}

// Run is a long running routine used to manage the task
func (t *TaskRunner) Run() {
	t.logger.Printf("[DEBUG] client: starting task context for '%s' (alloc '%s')", t.task.Name, t.ctx.Alloc().ID)

	// TODO: Start
	for {
		select {
		case update := <-t.updateCh:
			// TODO: Update
			t.task = update
		case <-t.destroyCh:
			// TODO: Destroy
			return
		}
	}
}

// Update is used to update the task of the context
func (t *TaskRunner) Update(update *structs.Task) {
	select {
	case t.updateCh <- update:
	default:
		t.logger.Printf("[ERR] client: dropping task update '%s' (alloc '%s')", update.Name, t.ctx.Alloc().ID)
	}
}

// Destroy is used to indicate that the task context should be destroyed
func (t *TaskRunner) Destroy() {
	t.destroyLock.Lock()
	defer t.destroyLock.Unlock()

	if t.destroy {
		return
	}
	t.destroy = true
	close(t.destroyCh)
}