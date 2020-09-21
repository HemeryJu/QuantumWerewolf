package main

import (
	"flag"

	"google.golang.org/grpc"

	"github.com/HemeryJu/QuantumWerewolf/pb"
	"github.com/HemeryJu/QuantumWerewolf/server"
)

var serverOptionsPath string

func init() {
	flag.StringVar(&serverOptionsPath, "server_options", "server.json", "config file of the grpc listener")
}

func main() {
	options := server.GetServerOptions(serverOptionsPath)
	svc := server.NewServer(options)

	listener := grpc.NewServer()

	pb.RegisterWerewolfServerServer(listener, svc)
}
