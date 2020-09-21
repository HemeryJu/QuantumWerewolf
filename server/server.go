package server

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/HemeryJu/QuantumWerewolf/pb"
)

type Options struct {
	Port int `json:"port"`
}

func GetServerOptions(path string) (*Options, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Options
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
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