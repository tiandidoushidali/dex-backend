package etherum

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type Task struct {
}

func NewEthereumTask() *Task {
	return &Task{}
}

func RegisterTask(c *cron.Cron) {
	task := NewEthereumTask()
	entryID, err := c.AddJob("@every 1s", task)
	if err != nil {
		panic(err)
	}
	logx.Infof("RegisterTask task entryID: %d", entryID)

	eventPairCreateTask := NewEventPairCreate()
	entryID, err = c.AddJob("@every 2s", eventPairCreateTask)
	if err != nil {
		panic(err)
	}
	logx.Infof("RegisterTask eventPairCreateTask entryID: %d", entryID)
}

func (ethereum *Task) Run() {
	fmt.Println("ethereum task run")
}
