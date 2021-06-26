package worker

import (
	"bufio"
	"os/exec"
)

func StartJob(name string, argv []string, env []string) (*job, error) {
	execpath, err := exec.LookPath(name)
	if err != nil {
		return nil, err
	}

	command := exec.Command(execpath, argv...)
	job, err := NewJob(command) // Scheduled

	if err != nil {
		return nil, err
	}

	command.Env = env
	writer := bufio.NewWriter(job.outputFile)
	command.Stderr = writer
	command.Stdout = writer

	err = command.Start()
	if err != nil {
		return nil, err
	}
	job.JobState = Running

	return job, nil
}
