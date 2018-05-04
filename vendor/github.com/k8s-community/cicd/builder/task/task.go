package task

const (
	// StatePending marks what task is waiting for a free worker or processing
	StatePending = "pending"

	// StateSuccess marks what task was processed successfully
	StateSuccess = "success"

	// StateFailure marks what task was failed (for example, tests were failed)
	StateFailure = "failure"

	// StateError marks what task wasn't processed because of error on CI/CD system
	StateError = "error"
)

const (
	// TypeTest represents test task (it runs 'make test' command)
	TypeTest = "test"

	// TypeBuild represents build task (it runs 'make test' and 'make build' commands)
	TypeBuild = "build"
)

// Callback is a function to update information about current task state
type Callback func(taskID string, status string, description string)

// CICD represents a task for CI/CD.
type CICD struct {
	Callback  Callback
	ID        string
	Type      string
	Prefix    string // Prefix represents a prefix part for GOPATH, e.g. github.com, gitlab.com
	Repo      string // Repo represent full a path to the repository, e.g. k8s-community/cicd
	Commit    string
	Version   string
	Namespace string
}

// NewCICD creates an instance of a task.
func NewCICD(callback Callback, id, taskType, prefix, repo, commit, version string, namespace string) *CICD {
	return &CICD{
		Callback:  callback,
		ID:        id,
		Type:      taskType,
		Prefix:    prefix,
		Repo:      repo,
		Commit:    commit,
		Version:   version,
		Namespace: namespace,
	}
}
