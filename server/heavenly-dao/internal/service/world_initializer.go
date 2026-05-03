package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cultivation-world/shared/types"
)

// RegionTemplate defines a region to be created during world initialization.
type RegionTemplate struct {
	ID               string
	Name             string
	ParentID         *string
	SpiritualDensity float64
	SpiritualTier    int
	DangerLevel      int
	Description      string
	Lore             string
	Resources        []ResourceTemplate
}

// ResourceTemplate defines a resource to spawn in a region.
type ResourceTemplate struct {
	Name        string
	Type        string
	Rarity      int
	BaseQuantity int
	RespawnRate float64
}

// SectTemplate defines a sect to be created during world initialization.
type SectTemplate struct {
	Name        string
	RegionID    string
	Alignment   string
	Description string
	MaxMembers  int
}

// NPCTemplate defines an NPC to spawn during world initialization.
type NPCTemplate struct {
	Name         string
	RegionID     string
	Realm        string
	SpiritualRoot string
	Role         string
	Personality  string
}

// WorldLoreEntry represents a piece of world history/lore.
type WorldLoreEntry struct {
	Title       string
	Description string
	Era         string
}

// WorldConfig holds all templates for world initialization.
type WorldConfig struct {
	Regions []RegionTemplate
	Sects   []SectTemplate
	NPCs    []NPCTemplate
	Lore    []WorldLoreEntry
}

// DefaultWorldConfig returns the default world configuration.
func DefaultWorldConfig() *WorldConfig {
	qingzhou := "qingzhou"
	yunzhou := "yunzhou"
	tianzhou := "tianzhou"
	jingzhou := "jingzhou"
	yongzhou := "yongzhou"
	jizhou := "jizhou"

	return &WorldConfig{
		Regions: []RegionTemplate{
			{
				ID: "tianzhou", Name: "天州", SpiritualDensity: 0.95, SpiritualTier: 9,
				DangerLevel: 8, Description: "中州大陆，灵气最为浓郁，强者云集",
				Lore: "传说中天道的核心所在，上古仙人大道崩碎之地",
				Resources: []ResourceTemplate{
					{Name: "九转金丹", Type: "pill", Rarity: 10, BaseQuantity: 2, RespawnRate: 0.01},
					{Name: "混沌灵石", Type: "stone", Rarity: 9, BaseQuantity: 5, RespawnRate: 0.05},
				},
			},
			{
				ID: "qingzhou", Name: "青州", SpiritualDensity: 0.70, SpiritualTier: 7,
				DangerLevel: 5, Description: "东方大州，灵气充沛，修仙者众多",
				Lore: "青云宗发源地，万年前曾有一位渡劫期大能在此飞升",
				Resources: []ResourceTemplate{
					{Name: "青灵草", Type: "herb", Rarity: 5, BaseQuantity: 50, RespawnRate: 0.3},
					{Name: "紫精铜矿", Type: "ore", Rarity: 6, BaseQuantity: 20, RespawnRate: 0.2},
				},
			},
			{
				ID: "yunzhou", Name: "云州", SpiritualDensity: 0.65, SpiritualTier: 6,
				DangerLevel: 4, Description: "云中仙境，终年云雾缭绕",
				Lore: "上古时期为仙家洞府，残留大量阵法遗迹",
				Resources: []ResourceTemplate{
					{Name: "云雾灵茶", Type: "herb", Rarity: 6, BaseQuantity: 30, RespawnRate: 0.4},
				},
			},
			{
				ID: "xuanzhou", Name: "玄州", SpiritualDensity: 0.60, SpiritualTier: 6,
				DangerLevel: 6, Description: "北方玄冰之地，极寒灵气",
				Lore: "玄冰宫所在，万年玄冰封印着上古妖兽",
				Resources: []ResourceTemplate{
					{Name: "玄冰晶", Type: "stone", Rarity: 7, BaseQuantity: 15, RespawnRate: 0.15},
				},
			},
			{
				ID: "jingzhou", Name: "荆州", SpiritualDensity: 0.55, SpiritualTier: 5,
				DangerLevel: 5, Description: "南方火灵之地，炼丹师聚集",
				Lore: "丹鼎宗祖地，传说中有上古丹炉埋藏于此",
				Resources: []ResourceTemplate{
					{Name: "火灵果", Type: "herb", Rarity: 5, BaseQuantity: 40, RespawnRate: 0.35},
					{Name: "赤炎石", Type: "ore", Rarity: 5, BaseQuantity: 25, RespawnRate: 0.25},
				},
			},
			{
				ID: "yangzhou", Name: "扬州", SpiritualDensity: 0.50, SpiritualTier: 5,
				DangerLevel: 4, Description: "繁华商道，交易兴盛",
				Lore: "万宝阁总阁所在地，天下财富汇聚之地",
				Resources: []ResourceTemplate{
					{Name: "灵石矿", Type: "stone", Rarity: 3, BaseQuantity: 100, RespawnRate: 0.5},
				},
			},
			{
				ID: "yongzhou", Name: "雍州", SpiritualDensity: 0.45, SpiritualTier: 4,
				DangerLevel: 6, Description: "西方戈壁荒漠，险恶之地",
				Lore: "血煞殿隐秘之处，埋葬着无数古修的遗骸",
				Resources: []ResourceTemplate{
					{Name: "血灵石", Type: "stone", Rarity: 7, BaseQuantity: 10, RespawnRate: 0.1},
				},
			},
			{
				ID: "jizhou", Name: "冀州", SpiritualDensity: 0.40, SpiritualTier: 4,
				DangerLevel: 5, Description: "中部平原，农耕发达",
				Lore: "凡人最为密集的州郡，偶有修仙者微服私访",
				Resources: []ResourceTemplate{
					{Name: "灵谷", Type: "herb", Rarity: 2, BaseQuantity: 200, RespawnRate: 0.6},
				},
			},
			{
				ID: "yuzhou", Name: "豫州", SpiritualDensity: 0.35, SpiritualTier: 3,
				DangerLevel: 3, Description: "中原腹地，历史悠久",
				Lore: "天机阁总部所在，天下情报汇聚之处",
				Resources: []ResourceTemplate{
					{Name: "天机符", Type: "talismans", Rarity: 6, BaseQuantity: 15, RespawnRate: 0.2},
				},
			},
			{
				ID: "bingzhou", Name: "冰州", SpiritualDensity: 0.75, SpiritualTier: 7,
				DangerLevel: 7, Description: "极北冰封之地，冰灵根修士的圣地",
				Lore: "传说中有上古冰仙遗留的洞府",
				ParentID: &qingzhou,
				Resources: []ResourceTemplate{
					{Name: "万年寒冰", Type: "stone", Rarity: 8, BaseQuantity: 5, RespawnRate: 0.05},
				},
			},
			{
				ID: "wuzhou", Name: "巫州", SpiritualDensity: 0.55, SpiritualTier: 5,
				DangerLevel: 8, Description: "南疆巫蛊之地，神秘莫测",
				Lore: "巫蛊之术的发源地，禁忌知识的封印处",
				ParentID: &jingzhou,
				Resources: []ResourceTemplate{
					{Name: "蛊虫卵", Type: "creature", Rarity: 7, BaseQuantity: 8, RespawnRate: 0.1},
				},
			},
			{
				ID: "qianzhou", Name: "黔州", SpiritualDensity: 0.40, SpiritualTier: 4,
				DangerLevel: 7, Description: "西南密林，瘴气弥漫",
				Lore: "妖兽最为密集的原始森林",
				ParentID: &yunzhou,
				Resources: []ResourceTemplate{
					{Name: "瘴气灵果", Type: "herb", Rarity: 6, BaseQuantity: 20, RespawnRate: 0.2},
				},
			},
			{
				ID: "liangzhou", Name: "凉州", SpiritualDensity: 0.50, SpiritualTier: 5,
				DangerLevel: 6, Description: "西北边陲，风沙漫天",
				Lore: "古战场遗迹，残魂游荡之地",
				ParentID: &yongzhou,
				Resources: []ResourceTemplate{
					{Name: "风灵砂", Type: "ore", Rarity: 5, BaseQuantity: 30, RespawnRate: 0.3},
				},
			},
			{
				ID: "yizhou", Name: "益州", SpiritualDensity: 0.60, SpiritualTier: 6,
				DangerLevel: 4, Description: "天府之国，灵脉纵横",
				Lore: "蜀山剑派所在地，剑修圣地",
				ParentID: &jizhou,
				Resources: []ResourceTemplate{
					{Name: "剑意石", Type: "stone", Rarity: 7, BaseQuantity: 10, RespawnRate: 0.15},
				},
			},
			{
				ID: "youzhou", Name: "幽州", SpiritualDensity: 0.65, SpiritualTier: 6,
				DangerLevel: 7, Description: "东北幽暗之地，鬼修聚集",
				Lore: "幽冥入口，生死交界之处",
				ParentID: &tianzhou,
				Resources: []ResourceTemplate{
					{Name: "幽冥花", Type: "herb", Rarity: 8, BaseQuantity: 7, RespawnRate: 0.08},
				},
			},
			{
				ID: "secret_realm_ancient", Name: "上古秘境", SpiritualDensity: 0.90, SpiritualTier: 9,
				DangerLevel: 10, Description: "上古仙人遗留的独立空间",
				Lore: "千年开启一次，内部灵气浓郁程度是天州的数倍",
				ParentID: &tianzhou,
				Resources: []ResourceTemplate{
					{Name: "仙灵草", Type: "herb", Rarity: 10, BaseQuantity: 3, RespawnRate: 0.01},
					{Name: "仙器碎片", Type: "artifact", Rarity: 10, BaseQuantity: 1, RespawnRate: 0.005},
				},
			},
		},
		Sects: []SectTemplate{
			{
				Name: "青云宗", RegionID: "qingzhou", Alignment: "righteous",
				Description: "正道大宗，以剑修闻名天下", MaxMembers: 5000,
			},
			{
				Name: "血煞殿", RegionID: "yongzhou", Alignment: "demonic",
				Description: "魔道巨擘，修炼血道功法", MaxMembers: 3000,
			},
			{
				Name: "天机阁", RegionID: "yuzhou", Alignment: "neutral",
				Description: "天下情报第一，神秘莫测", MaxMembers: 500,
			},
			{
				Name: "丹鼎宗", RegionID: "jingzhou", Alignment: "righteous",
				Description: "炼丹第一大宗", MaxMembers: 2000,
			},
			{
				Name: "玄冰宫", RegionID: "xuanzhou", Alignment: "righteous",
				Description: "修炼冰系功法的圣地", MaxMembers: 1000,
			},
			{
				Name: "蜀山剑派", RegionID: "yizhou", Alignment: "righteous",
				Description: "剑修圣地，天下剑修向往之所", MaxMembers: 3000,
			},
		},
		NPCs: []NPCTemplate{
			{Name: "青云子", RegionID: "qingzhou", Realm: "nascent_soul", SpiritualRoot: "剑", Role: "掌门", Personality: "威严"},
			{Name: "林清雪", RegionID: "qingzhou", Realm: "golden_core", SpiritualRoot: "冰", Role: "长老", Personality: "清冷"},
			{Name: "张小凡", RegionID: "qingzhou", Realm: "qi_condensation", SpiritualRoot: "混元", Role: "弟子", Personality: "坚毅"},
			{Name: "苏灵儿", RegionID: "qingzhou", Realm: "foundation", SpiritualRoot: "灵", Role: "内门弟子", Personality: "活泼"},
			{Name: "赵无极", RegionID: "qingzhou", Realm: "golden_core", SpiritualRoot: "雷", Role: "执法长老", Personality: "刚正"},
			{Name: "陈青云", RegionID: "qingzhou", Realm: "foundation", SpiritualRoot: "风", Role: "外门弟子", Personality: "豪爽"},
			{Name: "王雨萱", RegionID: "qingzhou", Realm: "qi_condensation", SpiritualRoot: "水", Role: "弟子", Personality: "温柔"},
			{Name: "刘天风", RegionID: "qingzhou", Realm: "golden_core", SpiritualRoot: "剑", Role: "首席弟子", Personality: "高傲"},
			{Name: "血魔老祖", RegionID: "yongzhou", Realm: "nascent_soul", SpiritualRoot: "血", Role: "殿主", Personality: "残忍"},
			{Name: "血玫瑰", RegionID: "yongzhou", Realm: "golden_core", SpiritualRoot: "火", Role: "圣女", Personality: "妖娆"},
			{Name: "厉无血", RegionID: "yongzhou", Realm: "foundation", SpiritualRoot: "金", Role: "护法", Personality: "冷酷"},
			{Name: "鬼面人", RegionID: "yongzhou", Realm: "golden_core", SpiritualRoot: "暗", Role: "暗卫", Personality: "阴险"},
			{Name: "天机老人", RegionID: "yuzhou", Realm: "soul_transformation", SpiritualRoot: "天机", Role: "阁主", Personality: "神秘"},
			{Name: "慕容晓", RegionID: "yuzhou", Realm: "golden_core", SpiritualRoot: "风", Role: "情报使", Personality: "机灵"},
			{Name: "诸葛明", RegionID: "yuzhou", Realm: "foundation", SpiritualRoot: "木", Role: "谋士", Personality: "睿智"},
			{Name: "司马青", RegionID: "yuzhou", Realm: "foundation", SpiritualRoot: "水", Role: "记录官", Personality: "严谨"},
			{Name: "丹阳真人", RegionID: "jingzhou", Realm: "nascent_soul", SpiritualRoot: "火", Role: "宗主", Personality: "和善"},
			{Name: "药灵儿", RegionID: "jingzhou", Realm: "golden_core", SpiritualRoot: "木", Role: "炼丹师", Personality: "善良"},
			{Name: "火云邪神", RegionID: "jingzhou", Realm: "foundation", SpiritualRoot: "火", Role: "护宗", Personality: "暴躁"},
			{Name: "灵药仙子", RegionID: "jingzhou", Realm: "golden_core", SpiritualRoot: "木", Role: "药园管事", Personality: "温和"},
			{Name: "冰魄仙子", RegionID: "xuanzhou", Realm: "nascent_soul", SpiritualRoot: "冰", Role: "宫主", Personality: "高冷"},
			{Name: "寒霜剑客", RegionID: "xuanzhou", Realm: "golden_core", SpiritualRoot: "冰", Role: "剑侍", Personality: "冷峻"},
			{Name: "雪无痕", RegionID: "xuanzhou", Realm: "foundation", SpiritualRoot: "水", Role: "弟子", Personality: "沉默"},
			{Name: "剑尘子", RegionID: "yizhou", Realm: "soul_transformation", SpiritualRoot: "剑", Role: "掌门", Personality: "超脱"},
			{Name: "李逍遥", RegionID: "yizhou", Realm: "golden_core", SpiritualRoot: "剑", Role: "大师兄", Personality: "潇洒"},
			{Name: "赵灵儿", RegionID: "yizhou", Realm: "foundation", SpiritualRoot: "灵", Role: "弟子", Personality: "灵动"},
			{Name: "酒剑仙", RegionID: "yizhou", Realm: "nascent_soul", SpiritualRoot: "剑", Role: "客卿长老", Personality: "洒脱"},
			{Name: "万宝楼主", RegionID: "yangzhou", Realm: "golden_core", SpiritualRoot: "金", Role: "商人", Personality: "精明"},
			{Name: "流浪散修", RegionID: "jizhou", Realm: "qi_condensation", SpiritualRoot: "土", Role: "散修", Personality: "随性"},
			{Name: "神秘老者", RegionID: "tianzhou", Realm: "soul_transformation", SpiritualRoot: "混元", Role: "隐士", Personality: "深不可测"},
			{Name: "妖兽王", RegionID: "qianzhou", Realm: "golden_core", SpiritualRoot: "妖", Role: "妖兽首领", Personality: "凶猛"},
			{Name: "巫蛊婆婆", RegionID: "wuzhou", Realm: "nascent_soul", SpiritualRoot: "蛊", Role: "巫蛊师", Personality: "诡谲"},
			{Name: "幽冥使者", RegionID: "youzhou", Realm: "golden_core", SpiritualRoot: "幽冥", Role: "鬼修", Personality: "阴森"},
			{Name: "凉州侠客", RegionID: "liangzhou", Realm: "foundation", SpiritualRoot: "风", Role: "侠客", Personality: "侠义"},
			{Name: "冰州猎户", RegionID: "bingzhou", Realm: "qi_condensation", SpiritualRoot: "冰", Role: "猎人", Personality: "粗犷"},
			{Name: "秘境守护者", RegionID: "secret_realm_ancient", Realm: "nascent_soul", SpiritualRoot: "仙", Role: "守护者", Personality: "庄严"},
			{Name: "白云散人", RegionID: "yunzhou", Realm: "foundation", SpiritualRoot: "风", Role: "散修", Personality: "淡泊"},
			{Name: "赤炎尊者", RegionID: "jingzhou", Realm: "golden_core", SpiritualRoot: "火", Role: "散修", Personality: "暴躁"},
			{Name: "青莲仙子", RegionID: "qingzhou", Realm: "foundation", SpiritualRoot: "木", Role: "散修", Personality: "清雅"},
			{Name: "紫电真君", RegionID: "tianzhou", Realm: "nascent_soul", SpiritualRoot: "雷", Role: "散修", Personality: "刚烈"},
			{Name: "黄沙客", RegionID: "liangzhou", Realm: "qi_condensation", SpiritualRoot: "土", Role: "佣兵", Personality: "豪放"},
			{Name: "碧水真人", RegionID: "jizhou", Realm: "golden_core", SpiritualRoot: "水", Role: "散修", Personality: "温和"},
			{Name: "金甲武士", RegionID: "tianzhou", Realm: "foundation", SpiritualRoot: "金", Role: "护卫", Personality: "忠诚"},
			{Name: "木灵童子", RegionID: "yizhou", Realm: "qi_condensation", SpiritualRoot: "木", Role: "药童", Personality: "纯真"},
			{Name: "暗影幽卫", RegionID: "youzhou", Realm: "foundation", SpiritualRoot: "暗", Role: "刺客", Personality: "冷酷"},
			{Name: "灵狐仙", RegionID: "qianzhou", Realm: "golden_core", SpiritualRoot: "妖", Role: "妖修", Personality: "妩媚"},
			{Name: "石巨人", RegionID: "bingzhou", Realm: "foundation", SpiritualRoot: "土", Role: "守卫", Personality: "迟钝"},
			{Name: "风语者", RegionID: "yunzhou", Realm: "qi_condensation", SpiritualRoot: "风", Role: "信使", Personality: "敏捷"},
			{Name: "血衣侯", RegionID: "yongzhou", Realm: "nascent_soul", SpiritualRoot: "血", Role: "诸侯", Personality: "霸道"},
			{Name: "天机侍女", RegionID: "yuzhou", Realm: "qi_condensation", SpiritualRoot: "水", Role: "侍女", Personality: "恭顺"},
			{Name: "丹霞散人", RegionID: "jingzhou", Realm: "foundation", SpiritualRoot: "火", Role: "炼丹师", Personality: "专注"},
		},
		Lore: []WorldLoreEntry{
			{Title: "创世神话", Description: "传说天道开辟天地，化万物为九州，灵气充盈其间", Era: "创世时代"},
			{Title: "第一次修仙大战", Description: "正道与魔道在青州决战，波及方圆万里", Era: "上古时代"},
			{Title: "天道法则确立", Description: "天机阁初代阁主观测天道，确立了修炼体系", Era: "上古时代"},
			{Title: "灵气衰退", Description: "末法时代来临，灵气浓度下降了七成", Era: "中古时代"},
			{Title: "新纪元开启", Description: "天道复苏，灵气开始回升，新的修炼时代到来", Era: "新纪元"},
		},
	}
}

// WorldInitializer handles the initialization of the world state.
type WorldInitializer struct {
	config  *WorldConfig
	regions map[types.RegionID]types.Region
	sects   map[string]*types.Sect
	npcs    []*types.Entity
	lore    []WorldLoreEntry
}

// NewWorldInitializer creates a new WorldInitializer with the given config.
func NewWorldInitializer(config *WorldConfig) *WorldInitializer {
	return &WorldInitializer{
		config:  config,
		regions: make(map[types.RegionID]types.Region),
		sects:   make(map[string]*types.Sect),
		npcs:    []*types.Entity{},
		lore:    []WorldLoreEntry{},
	}
}

// Initialize performs the complete world initialization.
func (w *WorldInitializer) Initialize(seed int64) (*WorldInitResult, error) {
	rng := rand.New(rand.NewSource(seed))

	if err := w.initializeRegions(rng); err != nil {
		return nil, fmt.Errorf("初始化区域失败: %w", err)
	}

	if err := w.initializeSects(rng); err != nil {
		return nil, fmt.Errorf("初始化宗门失败: %w", err)
	}

	if err := w.initializeNPCs(rng); err != nil {
		return nil, fmt.Errorf("初始化NPC失败: %w", err)
	}

	w.initializeLore()

	return &WorldInitResult{
		Regions: w.regions,
		Sects:   w.sects,
		NPCs:    w.npcs,
		Lore:    w.lore,
	}, nil
}

// WorldInitResult holds the result of world initialization.
type WorldInitResult struct {
	Regions map[types.RegionID]types.Region
	Sects   map[string]*types.Sect
	NPCs    []*types.Entity
	Lore    []WorldLoreEntry
}

func (w *WorldInitializer) initializeRegions(rng *rand.Rand) error {
	for _, tmpl := range w.config.Regions {
		id := types.RegionID(tmpl.ID)

		var parentID *types.RegionID
		if tmpl.ParentID != nil {
			pID := types.RegionID(*tmpl.ParentID)
			parentID = &pID
		}

		resources := make([]types.Resource, 0, len(tmpl.Resources))
		for _, rt := range tmpl.Resources {
			quantity := rt.BaseQuantity
			if quantity > 0 {
				variance := int(float64(quantity) * 0.2)
				if variance > 0 {
					quantity += rng.Intn(variance*2) - variance
				}
				if quantity < 1 {
					quantity = 1
				}
			}

			resources = append(resources, types.Resource{
				ID:          fmt.Sprintf("res_%s_%s", tmpl.ID, rt.Name),
				Name:        rt.Name,
				Type:        rt.Type,
				Rarity:      rt.Rarity,
				Quantity:    quantity,
				RespawnRate: rt.RespawnRate,
			})
		}

		region := types.Region{
			ID:               id,
			Name:             tmpl.Name,
			ParentRegionID:   parentID,
			SpiritualDensity: tmpl.SpiritualDensity,
			SpiritualTier:    tmpl.SpiritualTier,
			DangerLevel:      tmpl.DangerLevel,
			Resources:        resources,
			Description:      tmpl.Description,
			Lore:             tmpl.Lore,
			Rules: types.RegionRules{
				TaxRate: 0.05,
			},
		}

		w.regions[id] = region
	}

	return nil
}

func (w *WorldInitializer) initializeSects(rng *rand.Rand) error {
	for _, tmpl := range w.config.Sects {
		regionID := types.RegionID(tmpl.RegionID)
		if _, exists := w.regions[regionID]; !exists {
			return fmt.Errorf("宗门 '%s' 的区域 '%s' 不存在", tmpl.Name, tmpl.RegionID)
		}

		sect := &types.Sect{
			ID:        fmt.Sprintf("sect_%s", tmpl.Name),
			Name:      tmpl.Name,
			Alignment: tmpl.Alignment,
			MemberCount: 0,
			Prestige:  100,
			Territory: []string{tmpl.RegionID},
			Wealth:    10000,
			CreatedAt: time.Now().Unix(),
		}

		w.sects[sect.ID] = sect
	}

	return nil
}

func (w *WorldInitializer) initializeNPCs(rng *rand.Rand) error {
	for _, tmpl := range w.config.NPCs {
		regionID := types.RegionID(tmpl.RegionID)
		if _, exists := w.regions[regionID]; !exists {
			return fmt.Errorf("NPC '%s' 的区域 '%s' 不存在", tmpl.Name, tmpl.RegionID)
		}

		realm := parseRealm(tmpl.Realm)

		entity := &types.Entity{
			ID:          types.EntityID(fmt.Sprintf("npc_%s", tmpl.Name)),
			EntityType:  types.EntityTypeNPC,
			Name:        tmpl.Name,
			Realm:       realm,
			Position: types.WorldPosition{
				RegionID: string(regionID),
			},
			Attributes: types.Attributes{
				SpiritualRoots: []types.SpiritualRoot{
					{Element: tmpl.SpiritualRoot, Purity: 50 + rng.Intn(50)},
				},
				RootPurity: 50 + rng.Intn(50),
				Age:        18 + rng.Intn(200),
			},
			Karma:      types.Karma{KarmaValue: 0, Merit: 0, KarmicDebt: 0},
			Status:     types.StatusNormal,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		w.npcs = append(w.npcs, entity)
	}

	return nil
}

func (w *WorldInitializer) initializeLore() {
	w.lore = append(w.lore, w.config.Lore...)
}

func parseRealm(realmStr string) types.CultivationRealm {
	switch realmStr {
	case "qi_condensation":
		return types.RealmQiCondensation
	case "foundation":
		return types.RealmFoundation
	case "golden_core":
		return types.RealmGoldenCore
	case "nascent_soul":
		return types.RealmNascentSoul
	case "soul_transformation":
		return types.RealmSoulTransform
	case "tribulation":
		return types.RealmTribulation
	default:
		return types.RealmQiCondensation
	}
}

// GetRegionCount returns the number of initialized regions.
func (w *WorldInitializer) GetRegionCount() int {
	return len(w.regions)
}

// GetSectCount returns the number of initialized sects.
func (w *WorldInitializer) GetSectCount() int {
	return len(w.sects)
}

// GetNPCCount returns the number of initialized NPCs.
func (w *WorldInitializer) GetNPCCount() int {
	return len(w.npcs)
}

// GetLoreCount returns the number of lore entries.
func (w *WorldInitializer) GetLoreCount() int {
	return len(w.lore)
}
