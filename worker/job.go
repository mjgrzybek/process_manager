package worker

import (
	"bufio"
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

type job struct {
	*exec.Cmd
	state JobState

	startedDate time.Time
	exitedDate time.Time

	// TBD: output synchronization needed not to mess stdout and stdin between flushes?
	outputFile     *os.File
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
		Cmd: command,
		state: Scheduled,
		startedDate: time.Now(),
		outputFile: outputFile,
	}

	command.Env = env
	writer := bufio.NewWriter(outputFile)
	command.Stderr = writer
	command.Stdout = writer

	err = command.Start()
	if err != nil {
		return nil, err
	}
	job.state = Running

	return job, nil
}


func (j *job) Tail() (*tail.Tail, error) {
	t, err := tail.TailFile(j.outputFile.Name(), tail.Config{Follow: true})
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



