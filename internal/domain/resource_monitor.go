package domain

import (
	"context"
	"go-storage/internal/config"
	"go-storage/pkg/errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type ResourceMonitor struct {
	config             *config.FileServer
	currentMemoryUsage int64
	activeUploads      int32
	mutex              sync.RWMutex

	uploadSemaphore chan struct{}

	circuitState CircuitBreakerState
	failures     int32
	lastFailTime time.Time

	bufferPool sync.Pool
}

type CircuitBreakerState int32

const (
	CircuitClosed CircuitBreakerState = iota
	CircuitOpen
	CircuitHalfOpen
)

func NewResourceMonitor(config *config.FileServer) *ResourceMonitor {
	rm := &ResourceMonitor{
		config:          config,
		uploadSemaphore: make(chan struct{}, config.MaxConcurrentUploads),
		circuitState:    CircuitClosed,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return make([]byte, config.BufferSize)
			},
		},
	}

	go rm.startMonitoring()

	return rm
}

func (rm *ResourceMonitor) CanAllocateMemory(size int64) bool {
	if size > rm.config.MaxMemoryPerRequest {
		return false
	}

	currentUsage := atomic.LoadInt64(&rm.currentMemoryUsage)
	if currentUsage+size > rm.config.MaxTotalMemoryForFiles {
		return false
	}

	if rm.isMemoryPressureHigh() {
		return false
	}

	return true
}

func (rm *ResourceMonitor) AllocateMemory(size int64) bool {
	if !rm.CanAllocateMemory(size) {
		return false
	}

	atomic.AddInt64(&rm.currentMemoryUsage, size)
	return true
}

func (rm *ResourceMonitor) ReleaseMemory(size int64) {
	atomic.AddInt64(&rm.currentMemoryUsage, -size)
}

func (rm *ResourceMonitor) AcquireUploadSlot(ctx context.Context) error {
	if !rm.canExecute() {
		return errors.BadRequest("upload circuit breaker is open")
	}

	select {
	case rm.uploadSemaphore <- struct{}{}:
		atomic.AddInt32(&rm.activeUploads, 1)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (rm *ResourceMonitor) ReleaseUploadSlot() {
	select {
	case <-rm.uploadSemaphore:
		atomic.AddInt32(&rm.activeUploads, -1)
	default:
		// Semaphore was already released
	}
}

func (rm *ResourceMonitor) GetBuffer() []byte {
	return rm.bufferPool.Get().([]byte)
}

func (rm *ResourceMonitor) ReturnBuffer(buf []byte) {
	rm.bufferPool.Put(buf)
}

func (rm *ResourceMonitor) RecordSuccess() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	atomic.StoreInt32(&rm.failures, 0)
	if atomic.LoadInt32((*int32)(&rm.circuitState)) == int32(CircuitHalfOpen) {
		atomic.StoreInt32((*int32)(&rm.circuitState), int32(CircuitClosed))
	}
}

func (rm *ResourceMonitor) RecordFailure() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	failures := atomic.AddInt32(&rm.failures, 1)
	rm.lastFailTime = time.Now()

	if failures >= int32(rm.config.MaxFailuresBeforeOpen) {
		atomic.StoreInt32((*int32)(&rm.circuitState), int32(CircuitOpen))
	}
}

func (rm *ResourceMonitor) GetResourceStats() *ResourceStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &ResourceStats{
		MemoryUsage: ResourceUsage{
			Current:     atomic.LoadInt64(&rm.currentMemoryUsage),
			Limit:       rm.config.MaxTotalMemoryForFiles,
			SystemUsed:  int64(memStats.Alloc),
			SystemTotal: int64(memStats.Sys),
		},
		ActiveUploads: int(atomic.LoadInt32(&rm.activeUploads)),
		MaxUploads:    rm.config.MaxConcurrentUploads,
		CircuitState:  CircuitBreakerState(atomic.LoadInt32((*int32)(&rm.circuitState))),
		Failures:      int(atomic.LoadInt32(&rm.failures)),
	}
}

type ResourceStats struct {
	MemoryUsage   ResourceUsage       `json:"memory_usage"`
	ActiveUploads int                 `json:"active_uploads"`
	MaxUploads    int                 `json:"max_uploads"`
	CircuitState  CircuitBreakerState `json:"circuit_state"`
	Failures      int                 `json:"failures"`
}

type ResourceUsage struct {
	Current     int64 `json:"current"`
	Limit       int64 `json:"limit"`
	SystemUsed  int64 `json:"system_used"`
	SystemTotal int64 `json:"system_total"`
}

func (rm *ResourceMonitor) canExecute() bool {
	state := CircuitBreakerState(atomic.LoadInt32((*int32)(&rm.circuitState)))

	switch state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		rm.mutex.RLock()
		shouldTryHalfOpen := time.Since(rm.lastFailTime) >= rm.config.CircuitBreakerTimeout
		rm.mutex.RUnlock()

		if shouldTryHalfOpen {
			if atomic.CompareAndSwapInt32((*int32)(&rm.circuitState), int32(CircuitOpen), int32(CircuitHalfOpen)) {
				return true
			}
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

func (rm *ResourceMonitor) isMemoryPressureHigh() bool {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	if memStats.Sys == 0 {
		return false
	}

	memoryPressure := float64(memStats.Alloc) / float64(memStats.Sys)
	return memoryPressure >= rm.config.MemoryPressureThreshold
}

func (rm *ResourceMonitor) isCPUPressureHigh() bool {
	numGoroutines := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	cpuPressure := float64(numGoroutines) / float64(numCPU*100)
	return cpuPressure >= rm.config.CPUPressureThreshold
}

func (rm *ResourceMonitor) startMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rm.performHealthCheck()
	}
}

func (rm *ResourceMonitor) performHealthCheck() {
	stats := rm.GetResourceStats()

	_ = stats

	rm.cleanupExpiredResources()
}

func (rm *ResourceMonitor) cleanupExpiredResources() {
	// Cleanup logic for expired chunked upload sessions
	// This would typically interact with a repository to clean up old sessions
}
