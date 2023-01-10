package app

import (
	"context"
	"fmt"
	"time"
)

type CronJob interface {
	// String 可选的任务文字描述
	String() string
	// Interval 任务执行间隔, 调度协程会在此时间间隔内循环触发 Do 任务, 任务的触发间隔不考虑任务的执行时间
	Interval() time.Duration
	// Do 定时任务, 每 Interval 个时间间隔都将触发此任务,即便上一个任务可能因超时未执行完毕.
	// 其中 Do 的执行耗时应 <= Interval 本身, 否则将视为超时, 超时将触发 WhenTimeout
	Do(ctx context.Context) error
	// WhenError 当 Do 执行失败时触发的回调, 若 Do 执行失败且超时, 则 WhenError 和 WhenTimeout
	// 将同时被执行
	WhenError()
	// WhenTimeout 当定时任务未在规定时间内执行完毕时触发的回调, 当上一次 Do 执行超时时, 此 WhenTimeout 将和
	// Do 同时执行, 即 Do 和 WhenError 在同一个由 WhenTimeout 创建的子协程内。
	WhenTimeout()
}

type CronJobFunc struct{}

// String 任务文字描述
func (c *CronJobFunc) String() string { return "CronJobFunc" }

// Interval 执行调度间隔
func (c *CronJobFunc) Interval() time.Duration { return 5 * time.Second }

// Do 定时任务
func (c *CronJobFunc) Do() func(ctx context.Context) error {
	fmt.Printf("%s Run at %s\n", c.String(), time.Now().String())
	return nil
}

// WhenTimeout 当任务超时时执行的回调
func (c *CronJobFunc) WhenTimeout() {
	fmt.Printf("%s Timeout at %s\n", c.String(), time.Now().String())
}

// WhenError 当 Do 执行失败时触发的回调
func (c *CronJobFunc) WhenError() {
	return
}

type Scheduler struct {
	job    CronJob
	pctx   context.Context
	ctx    context.Context
	ticker *time.Ticker
	cancel context.CancelFunc
}

func (s *Scheduler) String() string { return s.job.String() }
func (s *Scheduler) Run()           { go s.Scheduler() }

// AtTime 到达任务的执行时间
func (s *Scheduler) AtTime() <-chan time.Time { return s.ticker.C }

// Do 执行任务
func (s *Scheduler) Do() {
	done := make(chan struct{}, 1)
	go func() {
		err := s.job.Do(s.ctx)
		done <- struct{}{} // 任务执行完毕

		if err != nil { // 此次任务执行发生错误
			s.job.WhenError()
		}
	}()

	select {
	case <-done:
		return
	case <-time.After(s.job.Interval()):
		// 单步任务执行时间超过了任务循环间隔,认为超时
		s.job.WhenTimeout()
	}
}

// Cancel 取消此定时任务
func (s *Scheduler) Cancel() {
	s.cancel()
	s.ticker.Stop()
}

// Scheduler 当时间到达时就启动一个任务协程
func (s *Scheduler) Scheduler() {
	for {
		// 每次循环都将创建一个新的 context.Context 避免超时情况下互相影响
		s.ctx, s.cancel = context.WithTimeout(s.pctx, s.job.Interval())
		select {
		case <-s.pctx.Done(): // 父节点被关闭,终止任务
			break
		case <-s.AtTime(): // 到达任务的执行时间, 创建一个新的事件任务
			go s.Do()
		}
	}
}
