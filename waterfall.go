package waterfall

import (
	"fmt"
	"github.com/ismdeep/log"
	"sync"
	"time"
)

const (
	StatusPending = 0
	StatusRunning = 1
	StatusDone    = 2
)

// funcInfo function information
type funcInfo struct {
	Name string
	Func func() error
}

// 信号
type signalInfo struct {
	name   string
	signal int
}

// waterfall struct of waterfall
type waterfall struct {
	FuncList     []funcInfo
	Dependencies map[string][]string
	status       int
	signals      chan signalInfo
}

// New create waterfall
func New() *waterfall {
	return &waterfall{
		FuncList:     []funcInfo{},
		Dependencies: make(map[string][]string),
		status:       StatusPending,
		signals:      make(chan signalInfo),
	}
}

// AddFunc add func
func (receiver *waterfall) AddFunc(name string, f func() error) {
	receiver.FuncList = append(receiver.FuncList, funcInfo{
		Name: name,
		Func: f,
	})
}

// AddDependency add dependency
func (receiver *waterfall) AddDependency(funcName string, dependencyFuncNames []string) {
	if receiver.Dependencies == nil {
		receiver.Dependencies = make(map[string][]string)
	}

	receiver.Dependencies[funcName] = dependencyFuncNames
}

// NaiveRun naive run
func (receiver *waterfall) NaiveRun() error {
	for _, f := range receiver.FuncList {
		startTime := time.Now().UnixMilli()
		log.Info("Run", "Name", f.Name, "Status", "started")
		if err := f.Func(); err != nil {
			log.Info("Run", "Name", f.Name, "Status", "failed")
			return err
		}
		endTime := time.Now().UnixMilli()
		log.Info("Run", "Name", f.Name, "Status", "done", "timeElapse(ms)", endTime-startTime)
	}

	return nil
}

func (receiver *waterfall) StartRunner() {
	for {
		signal := <-receiver.signals
		fmt.Println(signal)
	}
}

// Run start flow
func (receiver *waterfall) Run() error {
	// 1. 初始化状态表

	wg := &sync.WaitGroup{}
	for _, f := range receiver.FuncList {
		wg.Add(1)
		go func(f *funcInfo) {
			startTime := time.Now().UnixMilli()
			log.Info("Run", "Name", f.Name, "Status", "started")
			if err := f.Func(); err != nil {
				log.Info("Run", "Name", f.Name, "Status", "failed") // TODO
			}
			endTime := time.Now().UnixMilli()
			log.Info("Run", "Name", f.Name, "Status", "done", "timeElapse(ms)", endTime-startTime)
			wg.Done()
		}(&f)
	}
	wg.Wait()

	return nil
}
