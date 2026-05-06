package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/cultivation-world/ai-scheduler/internal/repository"
	"github.com/cultivation-world/ai-scheduler/internal/service"
	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	var npcRepo *repository.NPCRepository
	db, err := repository.NewDatabase(&cfg.Database)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v (NPC persistence disabled)", err)
	} else {
		defer db.Close()
		npcRepo = repository.NewNPCRepository(db)
		log.Println("Connected to PostgreSQL for NPC persistence")
	}

	gameConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Warning: Failed to connect to game-server: %v (NPC world interaction disabled)", err)
		gameConn = nil
	} else {
		defer gameConn.Close()
		log.Printf("Connected to game-server at %s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	}

	var gameClient service.GameServiceClient
	if gameConn != nil {
		gameClient = service.NewGameGrpcClient(gameConn)
	}

	aiSvc := service.NewAISchedulerService(cfg, gameClient)

	if npcRepo != nil {
		loadPersistentNPCs(aiSvc, npcRepo)
	}

	// Start NPC autonomous behavior loop
	go npcBehaviorLoop(aiSvc, npcRepo, cfg)

	grpcServer := grpc.NewServer()
	cultivation.RegisterAISchedulerServiceServer(grpcServer, aiSvc)

	port := 50052
	if p := os.Getenv("GRPC_PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("AI Scheduler Service starting on :%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// npcBehaviorLoop runs as a background goroutine, periodically driving NPC decisions.
func npcBehaviorLoop(aiSvc *service.AISchedulerService, npcRepo *repository.NPCRepository, cfg *config.Config) {
	// Stagger the first tick to avoid startup spike
	time.Sleep(time.Duration(5+rand.Intn(10)) * time.Second)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		npcIDs := aiSvc.GetNPCIDs()
		if len(npcIDs) == 0 {
			continue
		}

		// Randomize tick interval: 30-60 seconds per batch
		ticker.Reset(time.Duration(30+rand.Intn(31)) * time.Second)

		log.Printf("NPC behavior tick: processing %d NPCs", len(npcIDs))

		for _, npcID := range npcIDs {
			processNPC(aiSvc, npcRepo, npcID)
		}

		// Persist NPC memory periodically
		if npcRepo != nil {
			persistNPCMemories(aiSvc, npcRepo)
		}
	}
}

func processNPC(aiSvc *service.AISchedulerService, npcRepo *repository.NPCRepository, npcID string) {
	profile := aiSvc.GetNPC(npcID)
	if profile == nil {
		return
	}

	available := getAvailableActions(profile)

	// Call the LLM/template decision system
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	decision, err := aiSvc.ScheduleDecision(ctx, &cultivation.DecisionRequest{
		NpcId:            npcID,
		Context:          fmt.Sprintf("realm=%s,region=%s,goal=%s", profile.Realm, profile.CurrentRegion, profile.CurrentGoal),
		AvailableActions: available,
	})
	if err != nil {
		log.Printf("NPC %s decision error: %v", npcID, err)
		return
	}

	log.Printf("NPC %s decision: %s (source: %s, reasoning: %s)",
		npcID, decision.Action, decision.Source, decision.Reasoning)

	// Execute the action through game-server
	result := aiSvc.ExecuteNPCAction(npcID, decision.Action, decision.Params)
	if result != nil {
		if result.Success {
			log.Printf("NPC %s action %s succeeded: %s", npcID, decision.Action, result.Message)
		} else {
			log.Printf("NPC %s action %s failed: %s", npcID, decision.Action, result.Message)
		}
	}

	// Update NPC profile in registry
	if profile.CurrentGoal == "" || rand.Float64() < 0.1 {
		newGoal := generateNPCMotivation(profile)
		aiSvc.UpdateNPCGoal(npcID, newGoal)
	}

	// Persist profile if repository available
	if npcRepo != nil {
		profileRow := &repository.NPCProfileRow{
			NPCID:           profile.NPCID,
			EntityID:        profile.EntityID,
			PersonalityType: profile.PersonalityType,
			MoralAlignment:  profile.MoralAlignment,
			AmbitionLevel:   profile.AmbitionLevel,
			RiskTolerance:   profile.RiskTolerance,
			BackgroundStory: profile.BackgroundStory,
			CurrentGoal:     profile.CurrentGoal,
			CurrentRegion:   profile.CurrentRegion,
			Realm:           profile.Realm,
			Status:          decision.Action,
		}
		if err := npcRepo.SaveProfile(context.Background(), profileRow); err != nil {
			log.Printf("Failed to persist NPC %s profile: %v", npcID, err)
		}
	}
}

// getAvailableActions returns context-appropriate actions for an NPC.
func getAvailableActions(profile *service.NPCProfile) []string {
	base := []string{"cultivate", "meditate", "explore"}

	if profile.RiskTolerance > 0.6 {
		base = append(base, "combat", "gather")
	} else {
		base = append(base, "gather")
	}

	if profile.AmbitionLevel > 60 {
		base = append(base, "craft", "trade")
	}

	return base
}

// generateNPCMotivation creates a new goal based on NPC personality.
func generateNPCMotivation(profile *service.NPCProfile) string {
	goals := map[string][]string{
		"aggressive": {
			"寻求强大的对手突破自我",
			"收集天材地宝提升修为",
			"探索危险区域寻找机缘",
		},
		"cautious": {
			"稳固当前境界打好基础",
			"收集修炼资源储备",
			"研究功法提升实力",
		},
		"scholarly": {
			"参悟大道法则",
			"钻研炼丹炼器之道",
			"游历四方增长见识",
		},
		"gregarious": {
			"结交四方道友",
			"寻找志同道合之人",
			"探索天地秘境",
		},
		"balanced": {
			"提升修为境界",
			"探索未知区域",
			"积累修炼资源",
		},
	}

	personality := profile.PersonalityType
	options, exists := goals[personality]
	if !exists {
		options = goals["balanced"]
	}

	return options[rand.Intn(len(options))]
}

// loadPersistentNPCs loads NPC profiles from database into memory.
func loadPersistentNPCs(aiSvc *service.AISchedulerService, npcRepo *repository.NPCRepository) {
	ctx := context.Background()
	profiles, err := npcRepo.GetAllActiveProfiles(ctx)
	if err != nil {
		log.Printf("Failed to load persistent NPCs: %v", err)
		return
	}

	for _, row := range profiles {
		aiSvc.RegisterNPC(ctx, &cultivation.NPCRegisterRequest{
			NpcId:           row.NPCID,
			PersonalityType: row.PersonalityType,
			MoralAlignment:  row.MoralAlignment,
			AmbitionLevel:   int32(row.AmbitionLevel),
			RiskTolerance:   row.RiskTolerance,
			BackgroundStory: row.BackgroundStory,
			CurrentGoal:     row.CurrentGoal,
		})
		// Restore additional profile data
		aiSvc.UpdateNPCState(row.NPCID, row.Realm, row.CurrentRegion, row.Status)
	}

	log.Printf("Loaded %d persistent NPC profiles", len(profiles))
}

// persistNPCMemories saves NPC memories to the database.
func persistNPCMemories(aiSvc *service.AISchedulerService, npcRepo *repository.NPCRepository) {
	npcIDs := aiSvc.GetNPCIDs()
	ctx := context.Background()

	for _, npcID := range npcIDs {
		store := aiSvc.GetMemoryStore(npcID)
		store.Consolidate()

		memories := store.ToPersistableMemories()
		memRows := make([]*repository.NPCMemoryRow, 0, len(memories))
		for _, mem := range memories {
			memRow := &repository.NPCMemoryRow{
				NPCID:             npcID,
				MemoryType:        mem.MemoryType,
				Content:           mem.Content,
				Importance:        mem.Importance,
				RelatedEntityID:   mem.RelatedEntityID,
				RelatedEntityName: mem.RelatedEntityName,
			}
			if mem.ExpiresAt != nil {
				memRow.ExpiresAt = mem.ExpiresAt
			}
			memRows = append(memRows, memRow)
		}
		if len(memRows) > 0 {
			if err := npcRepo.SaveMemoriesBatch(ctx, npcID, memRows); err != nil {
				log.Printf("Failed to persist memories for NPC %s: %v", npcID, err)
			}
		}

		// Persist relationships
		for _, rel := range store.GetAllRelationships() {
			relRow := &repository.NPCRelationshipRow{
				NPCID:             npcID,
				TargetID:          rel.TargetID,
				TargetName:        rel.TargetName,
				RelationshipType:  rel.RelationshipType,
				Affinity:          rel.Affinity,
				Familiarity:       rel.Familiarity,
				InteractionCount:  rel.InteractionCount,
			}
			if err := npcRepo.UpsertRelationship(ctx, relRow); err != nil {
				log.Printf("Failed to persist relationship for NPC %s: %v", npcID, err)
			}
		}
	}
}
