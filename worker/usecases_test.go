package worker

import (
	"fmt"
	"testing"
	"time"
)


func TestUseCases(t *testing.T) {
	lsJob, err := StartJob("ls", nil, nil)
	if err != nil {
		t.Fatal(err)
	}


	pingJob, err := StartJob("ping", []string{"localhost", "-c10"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	sigtermIgnorerJob, err := StartJob("../tools/sigterm-ignorer/sigterm-ignorer.sh", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2*time.Second)

	type args struct {
		job *job
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "ls",
			args: args{
				job: lsJob,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "ping",
			args: args{
				job: pingJob,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "sigterm-ignorer",
			args: args{
				job: sigtermIgnorerJob,
			},
			want:    -1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := tt.args.job

			go func() {
				tail, err := job.OutputReader()
				if err != nil {
					t.Error(err)
				}
				for line := range tail.Lines {
					fmt.Println(line.Text)
				}
			}()

			time.Sleep(2 * time.Second)

			Stop(job)
			exitcode := job.ProcessState.ExitCode()


			if (err != nil) != tt.wantErr {
				t.Errorf("exitcode error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if exitcode != tt.want {
				t.Errorf("exitcode got = %v, want %v", exitcode, tt.want)
			}
		})
	}
}

