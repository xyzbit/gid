package core

import (
	"log"
	"time"
)

// Generator.
type IDGenerator struct {
	name   string
	cfg    *Config
	status *Status
	// 用于控制异步加载
	watchReStartInterval time.Duration
	watchRetryInterval   time.Duration
	ch                   chan struct{}
	// 当前批次已生成的数量
	num int32
	// 最后生成的id
	lastID int64
	loader StatusLoader
}

func NewGenerator(name string, sl StatusLoader, opts ...Option) *IDGenerator {
	cfg := &Config{
		LoadNum:       5000,
		LoadThreshold: 2000,
		IncrFirstID:   1,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	g := &IDGenerator{
		name:                 name,
		cfg:                  cfg,
		watchReStartInterval: WatchReStartInterval,
		watchRetryInterval:   WatchRetryInterval,
		ch:                   make(chan struct{}, 1),
		loader:               sl,
	}
	if err := g.loadAndUpdateStatus(); err != nil {
		panic(err)
	}
	go g.watchLoadThreshold()
	return g
}

func (g *IDGenerator) NextID() int64 {
	g.num++
	g.lastID++
	if g.lastID > g.status.LastID {
		return -1 // 表示已经用完了
	}

	if (g.cfg.LoadNum - g.num) == g.cfg.LoadThreshold {
		g.ch <- struct{}{}
	} else if g.num == g.cfg.LoadNum {
		g.num = 0
	}
	return g.lastID
}

func (g *IDGenerator) loadAndUpdateStatus() error {
	currStatus, err := g.loader.Status(g.name)
	if err != nil {
		return err
	}
	if currStatus.ModuleID == 0 { // 表示第一次创建
		currStatus.LastID = g.cfg.IncrFirstID
	}
	g.status = currStatus
	g.lastID = currStatus.LastID
	g.status.LastID += int64(g.cfg.LoadNum)
	return g.loader.SaveStatus(g.name, g.status)
}

// watchLoadThreshold 监听阀值，当阀值达到时，异步加载直到成功.
func (g *IDGenerator) watchLoadThreshold() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error, watch load threshold error: %+v", err)
			<-time.After(WatchReStartInterval * time.Second)
			g.watchLoadThreshold()
		}
	}()
	for {
		_, ok := <-g.ch
		if !ok {
			return
		}
		g.status.LastID += int64(g.cfg.LoadNum)
		if err := g.loader.SaveStatus(g.name, g.status); err != nil {
			// 存在错误，异步加载直到成功
			log.Printf("Error, update status error: %+v", err)
			<-time.After(WatchRetryInterval)
			g.ch <- struct{}{}
		}
	}
}
