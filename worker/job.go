package worker

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type JobState int
const (
	Scheduled JobState = iota
	Running = iota
	Stopped = iota
)

func (js JobState) String() string {
	return [...]string{"Scheduled", "Running", "Stopped"}[js]
}

type jobStatus struct {
	JobState
}

type job struct {
	*exec.Cmd
	jobStatus

	// TBD: output synchronization needed not to mess stdout and stdin between flushes?
	combinedOutputFile *os.File
}

func NewJob(command *exec.Cmd) (*job, error) {
	logsStorageDir := path.Join(os.TempDir(), "process_runner")

	if _, err := os.Stat(logsStorageDir); os.IsNotExist(err) {
		err = os.Mkdir(logsStorageDir, 0744)
		if err != nil {
			return nil, err
		}
	}

	outputFile, err := ioutil.TempFile(logsStorageDir, "output_")
	if err != nil {
		return nil, err
	}

	return &job{
		Cmd: command,
		jobStatus: jobStatus{
			JobState: Scheduled,
		},
		combinedOutputFile: outputFile,
	}, nil
}
