package api

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/leapp-to/leapp-go/pkg/executor"
)

type ActorRunnerID struct {
	uuid.UUID
}

func NewActorRunnerID() ActorRunnerID {
	u, _ := uuid.NewRandom()
	return ActorRunnerID{u}
}

type ActorRunner struct {
	registry map[ActorRunnerID]*taskEntry
	create   chan *actorCreateParams
	status   chan *actorStatusParams
	stream   chan *actorStreamParams
}

type actorStreamParams struct {
	ID       ActorRunnerID
	Target   io.Writer
	Response chan error
}

type actorStatusParams struct {
	ID       ActorRunnerID
	Response chan *ActorStatus
}

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

func (a *ActorRunner) StreamOutput(id ActorRunnerID, target io.Writer) error {
	result := make(chan error)
	a.stream <- &actorStreamParams{
		ID:       id,
		Target:   target,
		Response: result,
	}
	return <-result
}

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

func (a *ActorRunner) GetStatus(id ActorRunnerID) *ActorStatus {
	params := actorStatusParams{
		ID:       id,
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

type NoSuchTaskError struct{}

func (n NoSuchTaskError) Error() string {
	return "The requested task has not been found"
}

func (a *ActorRunner) doStream(params *actorStreamParams) {
	if entry, ok := a.registry[params.ID]; !ok {
		params.Response <- NoSuchTaskError{}
	} else {
		file, err := os.Open(entry.StreamPath)
		params.Response <- err
		if file != nil {

			go func(w io.Writer, f io.ReadCloser, doneChan chan struct{}) {
				defer f.Close()
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
			}(params.Target, file, entry.Done)
		}
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
