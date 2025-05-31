package solana

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type Task struct {
}

func NewSolanaTask() *Task {
	return &Task{}
}

func RegisterTask(c *cron.Cron) {
	task := NewSolanaTask()
	entryID, err := c.AddJob("@every 1s", task)
	if err != nil {
		panic(err)
	}
	logx.Infof("RegisterTask entryID: %d", entryID)
}

func (t *Task) Run() {
	fmt.Println("solana task run")
}
