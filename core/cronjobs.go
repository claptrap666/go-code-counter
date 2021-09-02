package core

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var counters map[uint]*cron.Cron = map[uint]*cron.Cron{}

func NewCounter() {
	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		Count()
	})
	c.Start()
}

func Count() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()
	go func(ctx context.Context) {
		//TODO: counter the code by commits
	}(ctx)
	select {
	case <-ctx.Done():
		logrus.Println("call successfully!!!")
		return
	case <-time.After(time.Duration(time.Hour * 2)):
		fmt.Println("timeout!!!")
		return
	}
}

func InitCounters() {
	NewCounter()
}
