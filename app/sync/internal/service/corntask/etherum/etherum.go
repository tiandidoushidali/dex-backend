package etherum

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

type Task struct {
}

func NewEtherumTask() *Task {
	return &Task{}
}

func RegisterTask(c *cron.Cron) {
	task := NewEtherumTask()
	c.AddJob("@every 1s", task)

	eventPairCreateTask := NewEventPairCreate()
	c.AddJob("@every 2s", eventPairCreateTask)
}

func (etherum *Task) Run() {
	fmt.Println("etherum task run")
}
