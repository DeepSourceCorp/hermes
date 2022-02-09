package infrastructure

import (
	"context"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

type Machinery struct {
	server *machinery.Server
}

var Server *Machinery

func GetMachineryServer(opts *config.Config) (*Machinery, error) {
	if Server != nil {
		return Server, nil
	}
	server, err := machinery.NewServer(opts)
	if err != nil {
		return nil, err
	}
	return &Machinery{server: server}, nil
}

func (im *Machinery) StartWorker(name string, concurrency int) error {
	worker := im.server.NewWorker(name, concurrency)
	if err := worker.Launch(); err != nil {
		return err
	}
	return nil
}

func (im *Machinery) RegisterTask(name string, task interface{}) error {
	return im.server.RegisterTask(name, task)
}

func (im *Machinery) SendTask(ctx context.Context, signature *tasks.Signature) error {
	_, err := im.server.SendTaskWithContext(ctx, signature)
	if err != nil {
		return err
	}
	return nil
}
