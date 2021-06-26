package worker

import (
	"log"
	"os"
	"syscall"
	"time"
)

func (j *job) Stop() {
	processState := make(chan *os.ProcessState)
	go processWait(j.Process, processState)

	_ = j.Process.Signal(syscall.SIGTERM)

	select {
	case j.ProcessState = <-processState:
	case <-time.After(3 * time.Second):
		// no one can resist SIGKILL
		_ = j.Process.Signal(syscall.SIGKILL)
		j.ProcessState = <-processState
	}

	j.exitedDate = time.Now()
	j.state = Stopped
}

func processWait(process *os.Process, state chan *os.ProcessState) {
	ps, err := process.Wait()
	// TODO: handle error
	if err != nil {
		log.Printf("ProcessWait: %e", err)
	}
	state <- ps
}
