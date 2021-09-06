package waterfall

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func LoadConfig() error {
	time.Sleep(time.Duration(rand.Intn(300)+1) * time.Millisecond)
	return nil
}

func LoadMonitor() error {
	time.Sleep(time.Duration(rand.Intn(300)+1) * time.Millisecond)
	return nil
}

func LoadMessenger() error {
	time.Sleep(time.Duration(rand.Intn(300)+1) * time.Millisecond)
	return nil
}

func LoadAFailedFunc() error {
	time.Sleep(time.Duration(rand.Intn(300)+1) * time.Millisecond)
	return errors.New("failed")
}

func TestNew(t *testing.T) {
	pool := New()
	assert.NotNil(t, pool.FuncList)
	assert.NotNil(t, pool.Dependencies)
}

func Test_waterfall_AddFunc(t *testing.T) {
	pool := New()
	pool.AddFunc("config", LoadConfig)
	pool.AddFunc("monitor", LoadMonitor)
	pool.AddFunc("messenger", LoadMessenger)
	assert.Equal(t, 3, len(pool.FuncList))
}

func Test_waterfall_AddDependency(t *testing.T) {
	pool := New()
	pool.AddFunc("config", LoadConfig)
	pool.AddFunc("monitor", LoadMonitor)
	pool.AddFunc("messenger", LoadMessenger)
	pool.AddDependency("monitor", []string{"config"})
	pool.AddDependency("messenger", []string{"config"})
	assert.Equal(t, 3, len(pool.FuncList))
	assert.Equal(t, pool.Dependencies["monitor"], []string{"config"})
	assert.Equal(t, pool.Dependencies["messenger"], []string{"config"})
}

func Test_waterfall_AddDependency2(t *testing.T) {
	pool := New()
	pool.Dependencies = nil
	pool.AddFunc("config", LoadConfig)
	pool.AddFunc("monitor", LoadMonitor)
	pool.AddFunc("messenger", LoadMessenger)
	pool.AddDependency("monitor", []string{"config"})
	pool.AddDependency("messenger", []string{"config"})
	assert.Equal(t, 3, len(pool.FuncList))
	assert.Equal(t, pool.Dependencies["monitor"], []string{"config"})
	assert.Equal(t, pool.Dependencies["messenger"], []string{"config"})
}

func Test_waterfall_Run(t *testing.T) {
	pool := New()
	pool.AddFunc("config", LoadConfig)
	pool.AddFunc("monitor", LoadMonitor)
	pool.AddFunc("messenger", LoadMessenger)
	pool.AddDependency("monitor", []string{"config"})
	pool.AddDependency("messenger", []string{"config"})
	pool.Run()
}

func Test_waterfall_Run_Failed(t *testing.T) {
	pool := New()
	pool.AddFunc("config", LoadConfig)
	pool.AddFunc("monitor", LoadMonitor)
	pool.AddFunc("messenger", LoadMessenger)
	pool.AddFunc("failed", LoadAFailedFunc)
	pool.AddDependency("monitor", []string{"config"})
	pool.AddDependency("messenger", []string{"config"})
	pool.Run()
}

func Test_waterfall_NaiveRun(t *testing.T) {
	pool := New()
	pool.AddFunc("config", LoadConfig)
	pool.AddFunc("monitor", LoadMonitor)
	pool.AddFunc("messenger", LoadMessenger)
	pool.AddDependency("monitor", []string{"config"})
	pool.AddDependency("messenger", []string{"config"})
	assert.NoError(t, pool.NaiveRun())
}

func Test_waterfall_NaiveRun_Failed(t *testing.T) {
	pool := New()
	pool.AddFunc("config", LoadConfig)
	pool.AddFunc("monitor", LoadMonitor)
	pool.AddFunc("messenger", LoadMessenger)
	pool.AddFunc("failed", LoadAFailedFunc)
	pool.AddDependency("monitor", []string{"config"})
	pool.AddDependency("messenger", []string{"config"})
	assert.Error(t, pool.NaiveRun())
}
