package solana

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

type Task struct {
}

func NewSolanaTask() *Task {
	return &Task{}
}

func RegisterTask(c *cron.Cron) {
	task := NewSolanaTask()
	c.AddJob("@every 1s", task)
}

func (etherum *Task) Run() {
	fmt.Println("solana task run")
}
