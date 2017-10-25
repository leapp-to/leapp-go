package api

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/leapp-to/leapp-go/pkg/executor"
)

// ActorRunnerID represents an unique of a running actor.
type ActorRunnerID struct {
	uuid.UUID
}

// NewActorRunnerID returns a new ActorRunnerID.
func NewActorRunnerID() ActorRunnerID {
	u, _ := uuid.NewRandom()
	return ActorRunnerID{u}
}

// ActorRunner stores running tasks along with channels to deliver results.
type ActorRunner struct {
	registry map[ActorRunnerID]*taskEntry
	create   chan *actorCreateParams
	status   chan *actorStatusParams
	stream   chan *actorStreamParams
}

type actorStreamParams struct {
	ID       ActorRunnerID
	Target   io.WriteCloser
	Response chan error
}

type actorStatusParams struct {
	ID       ActorRunnerID
	Response chan *ActorStatus
}

// ActorStatus represents the status of a running actor.
type ActorStatus struct {
	ID     ActorRunnerID
	Result *executor.Result
}

type actorCreateParams struct {
	ActorName string
	Stdin     string
	Response  chan ActorRunnerID
}

type taskEntry struct {
	ID         ActorRunnerID
	Done       chan struct{}
	StreamPath string
	Result     *executor.Result
	Error      error
}

// NewActorRunner returns a new ActorRunner.
func NewActorRunner() *ActorRunner {
	runner := &ActorRunner{
		registry: make(map[ActorRunnerID]*taskEntry),
		create:   make(chan *actorCreateParams),
		status:   make(chan *actorStatusParams),
		stream:   make(chan *actorStreamParams),
	}
	go runner.run()
	return runner
}

// StreamOutput sends a new actirStreanParams struct to an appropriate channel and returns an error.
func (a *ActorRunner) StreamOutput(id ActorRunnerID, target io.WriteCloser) error {
	result := make(chan error)
	a.stream <- &actorStreamParams{
		ID:       id,
		Target:   target,
		Response: result,
	}
	return <-result
}

// Create sends a new actirCreateParams struct to an appropriate channel to be executed and returns a pointer to a new ActorRunnerID.
func (a *ActorRunner) Create(actorName, stdin string) *ActorRunnerID {
	params := actorCreateParams{
		ActorName: actorName,
		Stdin:     stdin,
		Response:  make(chan ActorRunnerID),
	}
	a.create <- &params
	id, ok := <-params.Response
	if ok {
		close(params.Response)
		return &id
	}

	return nil
}

// GetStatus returns a status of a given ActorRunnerID.
func (a *ActorRunner) GetStatus(id *ActorRunnerID) *ActorStatus {
	params := actorStatusParams{
		ID:       *id,
		Response: make(chan *ActorStatus),
	}
	a.status <- &params
	status, ok := <-params.Response
	if ok {
		close(params.Response)
	}

	return status

}

func (a *ActorRunner) run() {
	for {
		select {
		case v, ok := <-a.create:
			if !ok {
				return
			}
			a.doCreate(v)
		case v, ok := <-a.status:
			if !ok {
				return
			}
			a.doStatus(v)
		case v, ok := <-a.stream:
			if !ok {
				return
			}
			a.doStream(v)
		}
	}
}

// NoSuchTaskError represents a customized error used when the task is not found.
type NoSuchTaskError struct{}

// Error implements the error interface for NoSuchTaskError.
func (n NoSuchTaskError) Error() string {
	return "The requested task has not been found"
}

func actorStream(w io.WriteCloser, f io.ReadCloser, doneChan chan struct{}) {
	defer f.Close()
	defer w.Close()
	if doneChan != nil {
	Loop:
		for {
			select {
			default:
				io.Copy(w, f)
				time.Sleep(time.Second)
			case <-doneChan:
				io.Copy(w, f)
				break Loop
			}
		}
	} else {
		io.Copy(w, f)
	}
}

func (a *ActorRunner) doStream(params *actorStreamParams) {
	if entry, ok := a.registry[params.ID]; !ok {
		// In case of errors we have to do the closing
		params.Target.Close()
		params.Response <- NoSuchTaskError{}
	} else {
		file, err := os.Open(entry.StreamPath)
		if file != nil {
			go actorStream(params.Target, file, entry.Done)
		} else {
			// In case of errors we have to do the closing
			params.Target.Close()
		}
		params.Response <- err
	}
}

func (a *ActorRunner) doCreate(params *actorCreateParams) {
	id := NewActorRunnerID()
	a.registry[id] = &taskEntry{
		ID:         id,
		StreamPath: filepath.Join("/tmp", id.String()),
		Done:       make(chan struct{}),
		Result:     nil,
		Error:      nil,
	}
	cmd := executor.New(params.ActorName, params.Stdin)
	cmd.StderrFile = a.registry[id].StreamPath

	go func(cmd *executor.Command, task *taskEntry) {
		task.Result, task.Error = cmd.Execute()
		close(task.Done)
	}(cmd, a.registry[id])

	params.Response <- id
}

func (a *ActorRunner) doStatus(params *actorStatusParams) {
	if entry, ok := a.registry[params.ID]; !ok {
		params.Response <- nil
	} else {
		params.Response <- &ActorStatus{
			ID:     params.ID,
			Result: entry.Result,
		}
	}
}
