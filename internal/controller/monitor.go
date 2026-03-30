package controller

import (
	"fabric-sdk/internal/service"
	"runtime"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

type MonitorController struct{}

var Monitor = new(MonitorController)

// Health 健康检查端点
func (c *MonitorController) Health(r *ghttp.Request) {
	status := "healthy"
	dbStatus := "up"
	
	// 检查数据库连接
	if err := service.PingDatabase(); err != nil {
		status = "unhealthy"
		dbStatus = "down"
	}
	
	health := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().Unix(),
		"checks": map[string]interface{}{
			"database": map[string]interface{}{
				"status": dbStatus,
			},
		},
		"system": map[string]interface{}{
			"goroutines": runtime.NumGoroutine(),
			"memory":     getMemoryStats(),
		},
	}
	
	r.Response.WriteJson(health)
}

// Metrics Prometheus指标端点
func (c *MonitorController) Metrics(r *ghttp.Request) {
	metrics := map[string]interface{}{
		"up":                       1,
		"go_goroutines":            runtime.NumGoroutine(),
		"go_memstats_alloc_bytes":  getMemoryAlloc(),
		"go_memstats_sys_bytes":    getMemorySys(),
		"database_connections_idle": getDatabaseIdleConns(),
		"database_connections_in_use": getDatabaseInUseConns(),
		"timestamp":                 time.Now().Unix(),
	}
	
	r.Response.WriteJson(metrics)
}

// PingDatabase 数据库ping检查
func (c *MonitorController) PingDatabase(r *ghttp.Request) {
	if err := service.PingDatabase(); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    503,
			Message: "Database unavailable",
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Database OK",
	})
}

// Ready 就绪检查
func (c *MonitorController) Ready(r *ghttp.Request) {
	r.Response.WriteJson(map[string]interface{}{
		"status": "ready",
	})
}

// 启动存活性检查
func (c *MonitorController) Live(r *ghttp.Request) {
	r.Response.WriteJson(map[string]interface{}{
		"status": "live",
	})
}

func getMemoryStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return map[string]interface{}{
		"alloc":      m.Alloc,
		"total_alloc": m.TotalAlloc,
		"sys":        m.Sys,
		"num_gc":     m.NumGC,
	}
}

func getMemoryAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func getMemorySys() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Sys
}

func getDatabaseIdleConns() int {
	if service.DB == nil {
		return 0
	}
	sqlDB, err := service.DB.DB()
	if err != nil {
		return 0
	}
	stats := sqlDB.Stats()
	return stats.Idle
}

func getDatabaseInUseConns() int {
	if service.DB == nil {
		return 0
	}
	sqlDB, err := service.DB.DB()
	if err != nil {
		return 0
	}
	stats := sqlDB.Stats()
	return stats.InUse
}
