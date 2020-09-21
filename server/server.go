package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/HemeryJu/QuantumWerewolf/pb"
)

type Options struct {
}

func GetServerOptions(path string) *Options {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("unable to read server config file, err: %s", err.Error()))
	}
	var config *Options
	err = json.Unmarshal(contents, config)
	if err != nil {
		panic(fmt.Sprintf("unable to parse server config file, err: %s", err.Error()))
	}
	return config
}

type Server struct {
	options *Options
}

func NewServer(options *Options) *Server {
	return &Server{
		options: options,
	}
}

func (s *Server) Start(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

var _ pb.WerewolfServerServer = &Server{}