package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mjgrzybek/process_manager/proto"
	"google.golang.org/grpc"
	"log"
)

var (
	//tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	uuid = flag.String("jobuuid", "", "Job uuid")
	//caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr = flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")
	//serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	tail := flag.Args()

	if len(tail) < 1 {
		log.Fatal("Command missing")
	}
	log.Println("tail: ", tail)

	method := tail[0]

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewProcessManagerServiceClient(conn)

	switch method {
	case "start":
		uuid, err := client.Start(context.TODO(), &proto.StartRequest{
			Name: tail[1],
			Args: tail[2:],
			Env:  nil,
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(uuid)
	case "stop":
		_, err = client.Stop(context.TODO(), &proto.StopRequest{Uuid: *uuid})
		if err != nil {
			log.Fatal(err)
		}
	case "status":
		status, err := client.Status(context.TODO(), &proto.StatusRequest{Uuid: *uuid})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(status)
	case "jobs":
		jobsRsp, err := client.Jobs(context.TODO(), &proto.JobsRequest{})
		if err != nil {
			log.Fatal(err)
		}
		jobs := jobsRsp.GetJobs()

		for _, job := range jobs {
			fmt.Println(job.String())
		}
	case "output":
	default:
		fmt.Println("Unknown command: " + method)
	}
}
