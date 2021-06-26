package worker

import (
	"log"
	"os"
	"syscall"
	"time"
)

func Stop(job *job) {
	processState := make(chan *os.ProcessState)
	go processWait(job.Process, processState)

	_ = job.Process.Signal(syscall.SIGTERM)

	select {
	case job.ProcessState = <-processState:
	case <-time.After(3 * time.Second):
		_ = job.Process.Signal(syscall.SIGKILL)
		job.ProcessState = <-processState
	}

	job.JobState = Stopped
}

func processWait(process *os.Process, state chan *os.ProcessState) {
	ps, err := process.Wait()
	// TODO: handle error
	if err != nil {
		log.Printf("ProcessWait: %e", err)
	}
	state <- ps
}
