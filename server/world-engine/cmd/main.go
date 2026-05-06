package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cultivation-world/shared/config"
	"github.com/cultivation-world/world-engine/internal/repository"
	"github.com/cultivation-world/world-engine/internal/service"
	cultivation "github.com/cultivation-world/shared/proto/pb"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	// 数据库持久化
	var worldRepo *repository.WorldRepository
	if cfg.Database.Host != "" {
		var err error
		worldRepo, err = repository.NewWorldRepository(
			cfg.Database.Host, cfg.Database.Port,
			cfg.Database.User, cfg.Database.Password, cfg.Database.Database,
		)
		if err != nil {
			log.Printf("Warning: failed to connect to database for world state persistence: %v", err)
		} else {
			defer worldRepo.Close()
			log.Println("World Engine: 已连接数据库，世界状态将持久化")
		}
	}

	worldSvc := service.NewWorldEngineService()

	// 从数据库加载持久化的世界状态
	if worldRepo != nil {
		epoch, metrics, err := worldRepo.LoadWorldState()
		if err != nil {
			log.Printf("Warning: failed to load world state: %v", err)
		} else if metrics != nil {
			worldSvc.RestoreState(epoch, metrics)
			log.Printf("World Engine: 已恢复世界状态，epoch=%d", epoch)
		}

		// 加载各区域资源数量
		for _, region := range worldSvc.GetAllRegions() {
			resources, err := worldRepo.LoadRegionResources(string(region.ID))
			if err != nil {
				log.Printf("Warning: failed to load resources for region %s: %v", region.ID, err)
				continue
			}
			worldSvc.RestoreRegionResources(region.ID, resources)
		}
		log.Println("World Engine: 已恢复区域资源数量")
	}

	// 资源刷新 ticker（每5分钟）
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("World Engine: 开始刷新资源...")
			worldSvc.AdvanceEpoch()
		}
	}()

	// 事件调度器（每30秒检查触发条件和事件生命周期）
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		var lastNotifiedEvents = make(map[string]bool)

		for range ticker.C {
			expired := worldSvc.UpdateActiveEvents()
			for _, e := range expired {
				log.Printf("事件结束: [%s] %s 在区域 %s", e.Type, e.Name, e.RegionID)
			}

			newEvents := worldSvc.TryTriggerEvents()
			for _, e := range newEvents {
				log.Printf("事件触发: [%s] %s 在区域 %s，持续到 %s",
					e.Type, e.Name, e.RegionID, e.EndTime.Format("15:04:05"))
				lastNotifiedEvents[e.ID] = true
			}

			allActive := worldSvc.GetAllActiveEvents()
			activeMap := make(map[string]bool)
			for _, e := range allActive {
				activeMap[e.ID] = true
			}
			for id := range lastNotifiedEvents {
				if !activeMap[id] {
					delete(lastNotifiedEvents, id)
				}
			}
		}
	}()

	// 世界状态持久化（每5分钟保存）
	if worldRepo != nil {
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				state := worldSvc.GetStateForPersistence()
				if err := worldRepo.SaveWorldState(state.Epoch, state.Metrics); err != nil {
					log.Printf("Warning: failed to save world state: %v", err)
				}
				if err := worldRepo.SaveRegionResources(worldSvc.GetAllRegionsMap()); err != nil {
					log.Printf("Warning: failed to save region resources: %v", err)
				}
			}
		}()
	}

	// 后台事件触发通知监控（用于实时广播）
	go func() {
		notifyChan := worldSvc.GetNotifyChan()
		for event := range notifyChan {
			log.Printf("实时通知: 事件 [%s] %s 已在区域 %s 触发",
				event.Type, event.Name, event.RegionID)
		}
	}()

	grpcServer := grpc.NewServer()
	cultivation.RegisterWorldServiceServer(grpcServer, worldSvc)

	port := 50054
	if p := os.Getenv("GRPC_PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("World Engine Service starting on :%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
