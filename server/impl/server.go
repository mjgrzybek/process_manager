package impl

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mjgrzybek/process_manager/proto"
	"github.com/mjgrzybek/process_manager/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

type UUID int

type server struct {
	m map[UUID]*worker.Job
	id UUID
	proto.UnimplementedProcessManagerServiceServer
}

func NewServer() *server {
	s := &server{
		m: make(map[UUID]*worker.Job),
	}
	return s
}

func (s *server) Jobs(_ context.Context, req *proto.JobsRequest) (rsp *proto.JobsResponse, err error) {
	log.Print("req" + req.String())

	jobs := make([]*proto.JobsResponse_Job, 1)

	for uuid, job := range s.m {
		jobs = append(jobs, &proto.JobsResponse_Job{
			Uuid:   string(uuid),
			Status: makeStatusResponse(job),
		})
	}

	rsp = &proto.JobsResponse{
		Jobs: jobs,
	}
	log.Print("rsp" + rsp.String())
	return
}

func (s *server) Start(_ context.Context, req *proto.StartRequest) (rsp *proto.StartResponse, err error) {
	log.Print("req" + req.String())

	job, err := worker.NewJob(req.GetName(), req.GetArgs(), nil)
	if err != nil {
		log.Print(err)
	}

	rsp = &proto.StartResponse{
		Uuid: strconv.Itoa(int(s.id)),
	}

	s.m[s.id] = job

	s.id++
	log.Print("rsp" + rsp.String())
	return
}
func (s *server) Stop(_ context.Context, req *proto.StopRequest) (rsp *proto.StopResponse, err error) {
	log.Print("req" + req.String())

	atoi, _ := strconv.Atoi(req.GetUuid())
	job := s.m[UUID(atoi)]

	job.Stop()
	rsp = &proto.StopResponse{}
	return
}
func (s *server) Status(_ context.Context, req *proto.StatusRequest) (rsp *proto.StatusResponse, err error) {
	log.Print("req" + req.String())

	atoi, _ := strconv.Atoi(req.GetUuid())
	job := s.m[UUID(atoi)]

	rsp = makeStatusResponse(job)

	log.Print("rsp" + rsp.String())
	return
}

func makeStatusResponse(job *worker.Job) (rsp *proto.StatusResponse) {
	rsp = &proto.StatusResponse{}

	date, _ := job.GetStartedDate()
	if job.GetState() == worker.Running {
		rsp.State = &proto.StatusResponse_StartedProcess_{
			StartedProcess: &proto.StatusResponse_StartedProcess{
				StartedDate: &timestamp.Timestamp{
					Seconds: date.Unix(),
					Nanos:   0,
				},
			},
		}
	} else {
		exitedDate, _ := job.GetExitedDate()

		code, _ := job.GetExitCode()
		systemStatus, _ := job.GetSystemStatus()

		rsp.State = &proto.StatusResponse_ExitedProcess_{
			ExitedProcess: &proto.StatusResponse_ExitedProcess{
				StartedDate: &timestamp.Timestamp{
					Seconds: date.Unix(),
					Nanos:   0,
				},
				ExitedDate: &timestamp.Timestamp{
					Seconds: exitedDate.Unix(),
					Nanos:   0,
				},
				Exitcode:     int32(code),
				SystemStatus: systemStatus,
			},
		}
	}

	return rsp
}
func (s *server) OutputStream(*proto.OutputRequest, proto.ProcessManagerService_OutputStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method OutputStream not implemented")
}
