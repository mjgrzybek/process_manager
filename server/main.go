package main

import (
	"flag"
	"fmt"
	"github.com/mjgrzybek/process_manager/proto"
	"github.com/mjgrzybek/process_manager/server/impl"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	//tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile   = flag.String("cert_file", "", "The TLS cert file")
	//keyFile    = flag.String("key_file", "", "The TLS key file")
	//jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterProcessManagerServiceServer(grpcServer, impl.NewServer())
	_ = grpcServer.Serve(lis)
}
