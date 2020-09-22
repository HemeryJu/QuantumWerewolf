package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/HemeryJu/QuantumWerewolf/pb"
	"github.com/HemeryJu/QuantumWerewolf/server"
)

var serverOptionsPath string

func init() {
	flag.StringVar(&serverOptionsPath, "server_options", "server.json", "config file of the grpc listener")
}

func main() {
	options, err := server.GetServerOptions(serverOptionsPath)
	if err != nil {
		panic(fmt.Sprintf("fail to fetch server config, err: %s", err.Error()))
	}
	svc := server.NewServer(options)

	grpcServer := grpc.NewServer()

	pb.RegisterWerewolfServerServer(grpcServer, svc)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", options.Port))
	if err != nil {
		panic(fmt.Sprintf("fail to create listener, err: %s", err.Error()))
	}

	err = grpcServer.Serve(lis)
	panic(fmt.Sprintf("fail to create listener, err: %s", err.Error()))
}
