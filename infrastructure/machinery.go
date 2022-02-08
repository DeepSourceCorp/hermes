package infrastructure

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

type InfraMachinery struct {
	server *machinery.Server
}

var server *InfraMachinery

func GetMachineryServer(opts *config.Config) (*InfraMachinery, error) {
	if server != nil {
		return server, nil
	}
	server, err := machinery.NewServer(opts)
	if err != nil {
		return nil, err
	}
	return &InfraMachinery{server: server}, nil
}

func (im *InfraMachinery) StartWorker(name string, concurrency int) error {
	worker := im.server.NewWorker(name, concurrency)
	if err := worker.Launch(); err != nil {
		return err
	}
	return nil
}
