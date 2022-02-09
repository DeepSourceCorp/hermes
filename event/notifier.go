package event

import (
	"context"
	"encoding/json"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/deepsourcelabs/hermes/infrastructure"
)

type Notifier interface {
	Dispatch(context.Context, *Event) error
}

type notifier struct {
	taskQueue *infrastructure.Machinery
}

func NewNotifier(taskQueue *infrastructure.Machinery) Notifier {
	return &notifier{
		taskQueue: taskQueue,
	}
}

func (n *notifier) Dispatch(ctx context.Context, event *Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	signature := &tasks.Signature{
		Name: "event-created",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: string(bytes),
			},
		},
	}
	err = n.taskQueue.SendTask(ctx, signature)
	if err != nil {
		return err
	}
	return nil
}
