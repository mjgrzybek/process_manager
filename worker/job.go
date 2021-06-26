package worker

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/hpcloud/tail"
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
	outputFile     *os.File
	outputFilePath string
}

func (j *job) Close() error {
	err := j.outputFile.Close()
	if err != nil {
		return err
	}
	return nil
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
		outputFilePath: outputFile.Name(),
		outputFile: outputFile,
	}, nil
}


func (j *job) OutputReader() (*tail.Tail, error) {
	t, err := tail.TailFile(j.outputFilePath, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			if j.ProcessState != nil { // process is still running
				t.Stop()
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	return t, nil
}