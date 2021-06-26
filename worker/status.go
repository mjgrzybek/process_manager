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

func (j *job) GetState() JobState {
	return j.state
}

func (j *job) GetExitCode() (int, error) {
	if j.state != Stopped {
		return -1, errors.New("Exitcode not available for running process")
	}
	return j.ProcessState.ExitCode(), nil
}

func (j *job) GetSystemStatus() (string, error) {
	if j.state != Stopped {
		return "", errors.New("SystemStatus not available for running process")
	}
	return j.ProcessState.String(), nil
}

func (j *job) GetStartedDate() (time.Time, error) {
	if j.state == Scheduled {
		return time.Time{}, errors.New("Process not started yet")
	}
	return j.startedDate, nil
}

func (j *job) GetExitedDate() (time.Time, error) {
	if j.state != Stopped {
		return time.Time{}, errors.New("Process is not stopped")
	}
	return j.exitedDate, nil
}
