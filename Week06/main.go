package main

/*
	参考 Hystrix 实现一个滑动窗口计数器。

	以上作业，要求提交到 GitHub 上面，Week06 作业提交地址：
	https://github.com/Go-000/Go-000/issues/81

	请务必按照示例格式进行提交，不要复制其他同学的格式，以免格式错误无法抓取作业。
*/

import (
	"sync"
	"time"
)

type bucket struct {
	val int64
}

type SlideWindow struct {
	windowSize int64

	bucketMutex *sync.RWMutex
	Buckets     map[int64]*bucket
	bucketCache []*bucket
}

// 窗口计数
func (s *SlideWindow) Reduce(now time.Time) int64 {
	s.bucketMutex.RLock()
	defer s.bucketMutex.RUnlock()
	var sum int64 = 0
	for t, bucket := range s.Buckets {
		if t > now.Unix()-s.windowSize {
			sum += bucket.val
		}
	}
	return sum
}

func (s *SlideWindow) removeBucket() {
	critical := time.Now().Unix() - s.windowSize
	for t, bucket := range s.Buckets {
		if t <= critical {
			s.bucketCache = append(s.bucketCache, bucket)
			delete(s.Buckets, t)
		}
	}
}

// 计数
func (s *SlideWindow) Inc() {
	s.Add(1)
}

// 计数
func (s *SlideWindow) Add(i int64) {
	s.bucketMutex.Lock()
	bucket := s.getCurrentBucket()
	bucket.val += i
	s.bucketMutex.Unlock()
	s.removeBucket()
}

// 返回滑动窗口最大bucket的统计数量
func (s *SlideWindow) Max(now time.Time) int64 {
	var max int64
	s.bucketMutex.RLock()
	defer s.bucketMutex.RUnlock()

	for timestamp, bucket := range s.Buckets {
		if timestamp >= now.Unix()-s.windowSize {
			if bucket.val > max {
				max = bucket.val
			}
		}
	}
	return max
}

// 计算 bucket 平均计数
func (s *SlideWindow) Avg(now time.Time) int64 {
	return s.Reduce(now) / s.windowSize
}

// 获取当前 bucket
func (s *SlideWindow) getCurrentBucket() *bucket {
	now := time.Now().Unix()
	if b, ok := s.Buckets[now]; ok {
		return b
	}

	if l := len(s.bucketCache); l > 0 {
		b := s.bucketCache[l-1]
		s.bucketCache = s.bucketCache[:l-1]
		return b
	}

	b := new(bucket)
	s.Buckets[now] = b
	return b
}

// 系统统计数据详情
func (s *SlideWindow) Stats() map[int64]*bucket {
	return s.Buckets
}

// 按秒级维度划分， 1秒 一个 bucket
// size 为 bucket 数量
func NewWindow(size int64) SlideWindow {
	if size <= 0 {
		panic("The size must be greater than 0")
	}
	return SlideWindow{
		windowSize:  size,
		bucketMutex: &sync.RWMutex{},
		Buckets:     make(map[int64]*bucket, size),
	}
}

// 统计请求中成功，失败，拒绝，超时等情况
// 在实际应用中，可统计更多等类型
type Collector struct {
	mu        *sync.RWMutex
	successes *SlideWindow
	failures  *SlideWindow
	rejects   *SlideWindow
	timeout   *SlideWindow
}

type Metric struct {
	Successes int64
	Failures  int64
	Rejects   int64
	Timeouts  int64
}

func (c Collector) Update(mtr Metric) {
	c.successes.Add(mtr.Successes)
	c.failures.Add(mtr.Failures)
	c.rejects.Add(mtr.Rejects)
	c.timeout.Add(mtr.Timeouts)
}
