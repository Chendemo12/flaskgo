package app

import (
	"context"
	"fmt"
	"time"
)

type CronJob interface {
	String() string                      // 可选的任务文字描述
	Do() func(ctx context.Context) error // 定时任务
	Interval() time.Duration             // 执行间隔
	WhenTimeout(ctx context.Context)     // 当定时任务未在规定时间内执行完毕时触发的回调
	WhenError(ctx context.Context)       // 当 Do 执行失败时触发的回调
}

type CronJobFunc struct{}

// String 任务文字描述
func (c *CronJobFunc) String() string { return "CronJobFunc" }

// Do 定时任务
func (c *CronJobFunc) Do() func(ctx context.Context) error {
	fmt.Printf("%s Run at %s\n", c.String(), time.Now().String())
	return nil
}

// WhenTimeout 当任务超时时执行的回调
func (c *CronJobFunc) WhenTimeout(ctx context.Context) {
	fmt.Printf("%s Timeout at %s\n", c.String(), time.Now().String())
}

// Interval 执行调度间隔
func (c *CronJobFunc) Interval() time.Duration {
	return 5 * time.Second
}

// WhenError 当 Do 执行失败时触发的回调
func (c *CronJobFunc) WhenError(ctx context.Context) {
	return
}

type Runner struct {
	job     CronJob
	context context.Context
	ticker  *time.Ticker
	cancel  context.CancelFunc
}

func (r *Runner) Run()                     { r.Scheduler() }
func (r *Runner) Cancel()                  { r.cancel() }
func (r *Runner) AtTime() <-chan time.Time { return r.ticker.C }
func (r *Runner) String() string           { return r.job.String() }

func (r *Runner) Do() {
}

func (r *Runner) Scheduler() {
	go func() {
		for {
			select {
			case <-r.context.Done():
				break
			case <-r.AtTime():
				go r.job.Do()
				// TODO: 当 job.Do() 超时时触发任务
			}
		}
	}()
}
