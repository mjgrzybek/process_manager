package worker

import (
	"fmt"
	"testing"
	"time"
)


func TestUseCases(t *testing.T) {
	lsJob, err := NewJob("ls", nil, nil)
	if err != nil {
		t.Fatal(err)
	}


	pingJob, err := NewJob("ping", []string{"localhost"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	sigtermIgnorerJob, err := NewJob("../tools/signal-ignorer/signal-ignorer.sh", []string{"SIGTERM"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// give time processes to start
	time.Sleep(1*time.Second)

	type args struct {
		job *job
	}
	tests := []struct {
		name     string
		args     args
		exitcode int
		status string
		wantErr  bool
	}{
		{
			name: "ls",
			args: args{
				job: lsJob,
			},
			exitcode: 0,
			status: "exit status 0",
			wantErr:  false,
		},
		{
			name: "ping",
			args: args{
				job: pingJob,
			},
			exitcode: -1,
			status: "signal: terminated",
			wantErr:  false,
		},
		{
			name: "sigterm-ignorer - SITGTERM ignored, SIGKILL should be used",
			args: args{
				job: sigtermIgnorerJob,
			},
			exitcode: -1,
			status: "signal: killed",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := tt.args.job

			tail, err := job.Tail()
			if err != nil {
				t.Error(err)
			}
			go func() {
				for line := range tail.Lines {
					fmt.Println(line.Text)
				}
			}()

			job.Stop()

			if (err != nil) != tt.wantErr {
				t.Errorf("exitcode error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if job.ProcessState.ExitCode() != tt.exitcode {
				t.Errorf("exitcode got = %v, exitcode %v", job.ProcessState.ExitCode(), tt.exitcode)
			}
			if job.ProcessState.String() != tt.status {
				t.Errorf("status got = %v, status %v", job.String(), tt.status)
			}
		})
	}
}
