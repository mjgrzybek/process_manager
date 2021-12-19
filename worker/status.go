package worker

import (
	"errors"
	"time"
)

type Status interface {
	GetState() JobState
	GetExitCode() (int, error)
	GetSystemStatus() (string, error)
	GetStartedDate() (time.Time, error)
	GetExitedDate() (time.Time, error)
}

func (j *Job) GetState() JobState {
	return j.state
}

func (j *Job) GetExitCode() (int, error) {
	j.RLock()
	defer j.RUnlock()

	if j.state != Stopped {
		return -1, errors.New("Exitcode not available for running process")
	}

	return j.ProcessState.ExitCode(), nil
}

func (j *Job) GetSystemStatus() (string, error) {
	j.RLock()
	defer j.RUnlock()

	if j.state != Stopped {
		return "", errors.New("SystemStatus not available for running process")
	}

	return j.ProcessState.String(), nil
}

func (j *Job) GetStartedDate() (time.Time, error) {
	j.RLock()
	defer j.RUnlock()

	return j.startedDate, nil
}

func (j *Job) GetExitedDate() (time.Time, error) {
	j.RLock()
	defer j.RUnlock()

	if j.state != Stopped {
		return time.Time{}, errors.New("Process is not stopped")
	}

	return j.exitedDate, nil
}
