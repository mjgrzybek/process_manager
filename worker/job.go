package worker

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

type JobState int
type CommonName string

const (
	Running JobState = iota
	Stopped          = iota
)

func (js JobState) String() string {
	return [...]string{"Running", "Stopped"}[js]
}

type job struct {
	*exec.Cmd
	state JobState

	startedDate time.Time
	exitedDate  time.Time

	outputFile *os.File

	sync.RWMutex
}

func (j *job) Close() error {
	err := j.outputFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewJob(name string, argv []string, env []string) (*job, error) {
	execpath, err := exec.LookPath(name)
	if err != nil {
		return nil, err
	}
	command := exec.Command(execpath, argv...)

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

	job := &job{
		Cmd:         command,
		state:       Running,
		startedDate: time.Now(),
		outputFile:  outputFile,
	}

	command.Env = env
	writer := bufio.NewWriter(outputFile)
	command.Stderr = writer
	command.Stdout = writer

	err = command.Start()
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (j *job) Tail() (*tail.Tail, error) {
	t, err := tail.TailFile(j.outputFile.Name(), tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			j.RLock()
			if j.ProcessState != nil { // process is still running
				_ = t.Stop()
			}
			j.RUnlock()

			time.Sleep(200 * time.Millisecond)
		}
	}()

	return t, nil
}
