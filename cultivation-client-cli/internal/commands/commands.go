package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"cultivation-client-cli/internal/client"

	"github.com/gorilla/websocket"
)

type cmdHandler func(conn *websocket.Conn, entityID string, args []string)

// CmdDef defines a command or command group.
type CmdDef struct {
	Names   []string   // ["角色", "role"] — first is display name
	Desc    string     // help description
	Usage   string     // usage template, "" = auto {name} {args}
	Handler cmdHandler // nil = group node
	Sub     []*CmdDef  // subcommands
	Context func() bool
}

// cmdTree is the root command list — groups + flat commands.
var cmdTree []*CmdDef

func init() {
	cmdTree = []*CmdDef{
		// ── 系统 ──
		{Names: []string{"帮助", "help"}, Desc: "显示帮助", Handler: cmdHelp},
		{Names: []string{"清屏", "clear", "cls"}, Desc: "清屏", Handler: cmdClear},
		{Names: []string{"退出", "exit", "quit"}, Desc: "退出游戏", Handler: cmdExit},

		// ── 角色 ──
		{Names: []string{"角色"}, Desc: "角色信息", Sub: []*CmdDef{
			{Names: []string{"状态", "status", "st"}, Desc: "查看角色状态", Handler: cmdStatus},
			{Names: []string{"属性", "attributes", "attrs"}, Desc: "查看详细属性", Handler: cmdAttributes},
			{Names: []string{"技能", "skills"}, Desc: "查看技能列表", Handler: cmdSkills},
				{Names: []string{"装备", "equip"}, Desc: "查看已装备物品", Handler: cmdShowEquipment},
				{Names: []string{"学习", "learn"}, Desc: "学习法术", Handler: cmdLearnSpell, Usage: "学习 <法术ID>"},
		}},

		// ── 修炼 ──
		{Names: []string{"修炼"}, Desc: "修炼相关", Sub: []*CmdDef{
			{Names: []string{"修炼", "cultivate", "cult"}, Desc: "运转功法修炼", Handler: wrapOp("cultivate")},
			{Names: []string{"打坐", "meditate", "med"}, Desc: "打坐恢复灵力神识", Handler: wrapOp("meditate")},
			{Names: []string{"休息", "sleep"}, Desc: "休息恢复满状态", Handler: wrapOp("sleep")},
			{Names: []string{"突破", "breakthrough", "bt"}, Desc: "突破境界", Handler: wrapOp("breakthrough")},
			{Names: []string{"探索", "explore", "exp"}, Desc: "探索当前区域", Handler: wrapOp("explore")},
			{Names: []string{"采集", "gather"}, Desc: "采集资源", Handler: cmdGather, Usage: "采集 <类型> [数量]"},
			{Names: []string{"移动", "move"}, Desc: "移动到区域", Handler: cmdMove, Usage: "移动 <区域ID>"},
			{Names: []string{"炼制", "craft"}, Desc: "炼制丹药/法器", Handler: cmdCraft, Usage: "炼制 <配方ID>"},
			{Names: []string{"自创", "create", "cm"}, Desc: "自创功法", Handler: wrapOp("create_method")},
		}},

		// ── 战斗 ──
		{Names: []string{"战斗"}, Desc: "战斗相关", Sub: []*CmdDef{
			{Names: []string{"目标", "attack", "atk"}, Desc: "攻击目标", Handler: cmdCombat, Usage: "目标 <目标ID>"},
			{Names: []string{"自动", "auto"}, Desc: "自动战斗模式", Handler: cmdAutoCombat},
			{Names: []string{"逃跑", "flee"}, Desc: "逃离战斗", Handler: wrapOp("flee"), Context: inCombat},
			{Names: []string{"技能", "skill"}, Desc: "使用技能/法术", Handler: cmdUseSkill, Usage: "技能 <技能ID> [目标ID]"},
		}},

		// ── 背包 ──
		{Names: []string{"背包", "bag", "inventory"}, Desc: "背包管理", Sub: []*CmdDef{
			{Names: []string{"列表", "list", "ls"}, Desc: "列出物品", Handler: cmdInventoryList},
			{Names: []string{"使用", "use"}, Desc: "使用物品", Handler: cmdUseItem, Usage: "使用 <物品名>"},
			{Names: []string{"丢弃", "drop"}, Desc: "丢弃物品", Handler: cmdDropItem, Usage: "丢弃 <物品ID>"},
		}},

		// ── 社交 ──
		{Names: []string{"社交", "social"}, Desc: "社交互动", Sub: []*CmdDef{
			{Names: []string{"聊天", "chat"}, Desc: "世界聊天", Handler: cmdChat, Usage: "聊天 <消息>"},
			{Names: []string{"私信", "message", "msg"}, Desc: "发送私信", Handler: cmdMessage, Usage: "私信 <收件人ID> <内容>"},
			{Names: []string{"好友", "friend"}, Desc: "好友管理", Sub: []*CmdDef{
				{Names: []string{"添加", "add"}, Desc: "添加好友", Handler: cmdAddFriend, Usage: "添加 <玩家名>"},
				{Names: []string{"删除", "remove", "del"}, Desc: "删除好友", Handler: cmdRemoveFriend, Usage: "删除 <好友ID>"},
				{Names: []string{"接受", "accept"}, Desc: "接受好友请求", Handler: cmdAcceptFriend, Usage: "接受 <请求ID>"},
					{Names: []string{"列表", "list", "ls"}, Desc: "查看好友列表", Handler: cmdListFriends},
			}},
		}},

		// ── 宗门 ──
		{Names: []string{"宗门", "sect"}, Desc: "宗门管理", Sub: []*CmdDef{
			{Names: []string{"创建", "create", "form"}, Desc: "创建宗门", Handler: cmdFormSect, Usage: "创建 <宗门名称>"},
			{Names: []string{"加入", "join"}, Desc: "加入宗门", Handler: cmdJoinSect, Usage: "加入 <宗门ID>"},
			{Names: []string{"退出", "leave"}, Desc: "退出宗门", Handler: cmdLeaveSect, Usage: "退出 <宗门ID>"},
				{Names: []string{"信息", "info"}, Desc: "查看宗门信息", Handler: cmdSectInfo, Usage: "信息 <宗门ID>"},
		}},

		// ── 交易 ──
		{Names: []string{"交易", "trade"}, Desc: "交易物品", Handler: cmdTrade, Usage: "交易 <目标ID> <物品ID> <价格>"},

		// ── 装备 ──
		{Names: []string{"装备", "equipment"}, Desc: "装备管理", Sub: []*CmdDef{
			{Names: []string{"列表", "list", "ls"}, Desc: "查看已装备物品", Handler: cmdShowEquipment},
			{Names: []string{"装备", "equip"}, Desc: "装备物品", Handler: cmdEquipItem, Usage: "装备 <物品ID> [装备位]"},
			{Names: []string{"卸下", "unequip"}, Desc: "卸下装备", Handler: cmdUnequipItem, Usage: "卸下 <装备位>"},
		}},

		// ── Flat aliases for backward compatibility (hidden from help) ──
		{Names: []string{"status", "st"}, Handler: cmdStatus, Context: func() bool { return false }},
		{Names: []string{"role"}, Handler: cmdStatus, Context: func() bool { return false }},
		{Names: []string{"cultivate", "cult"}, Handler: wrapOp("cultivate"), Context: func() bool { return false }},
		{Names: []string{"meditate", "med"}, Handler: wrapOp("meditate"), Context: func() bool { return false }},
		{Names: []string{"sleep"}, Handler: wrapOp("sleep"), Context: func() bool { return false }},
		{Names: []string{"breakthrough", "bt"}, Handler: wrapOp("breakthrough"), Context: func() bool { return false }},
		{Names: []string{"explore", "exp"}, Handler: wrapOp("explore"), Context: func() bool { return false }},
		{Names: []string{"gather"}, Handler: cmdGather, Context: func() bool { return false }},
		{Names: []string{"move"}, Handler: cmdMove, Context: func() bool { return false }},
		{Names: []string{"craft"}, Handler: cmdCraft, Context: func() bool { return false }},
		{Names: []string{"create_method", "cm"}, Handler: wrapOp("create_method"), Context: func() bool { return false }},
		{Names: []string{"combat", "fight"}, Handler: cmdCombat, Context: func() bool { return false }},
		{Names: []string{"flee"}, Handler: wrapOp("flee"), Context: func() bool { return false }},
		{Names: []string{"use_skill", "skill"}, Handler: cmdUseSkill, Context: func() bool { return false }},
		{Names: []string{"cast_spell", "spell"}, Handler: cmdUseSkill, Context: func() bool { return false }},
		{Names: []string{"chat"}, Handler: cmdChat, Context: func() bool { return false }},
		{Names: []string{"msg", "message", "send_message"}, Handler: cmdMessage, Context: func() bool { return false }},
		{Names: []string{"add_friend", "friend"}, Handler: cmdAddFriend, Context: func() bool { return false }},
		{Names: []string{"remove_friend", "unfriend"}, Handler: cmdRemoveFriend, Context: func() bool { return false }},
		{Names: []string{"accept_friend", "accept"}, Handler: cmdAcceptFriend, Context: func() bool { return false }},
		{Names: []string{"form_sect", "create_sect"}, Handler: cmdFormSect, Context: func() bool { return false }},
		{Names: []string{"join_sect"}, Handler: cmdJoinSect, Context: func() bool { return false }},
		{Names: []string{"leave_sect"}, Handler: cmdLeaveSect, Context: func() bool { return false }},
		{Names: []string{"trade"}, Handler: cmdTrade, Context: func() bool { return false }},
		{Names: []string{"learn_spell", "learn"}, Handler: cmdLearnSpell, Context: func() bool { return false }},
		{Names: []string{"list_friends"}, Handler: cmdListFriends, Context: func() bool { return false }},
		{Names: []string{"sect_info"}, Handler: cmdSectInfo, Context: func() bool { return false }},
		{Names: []string{"equip"}, Handler: cmdShowEquipment, Context: func() bool { return false }},
	}
}

func inCombat() bool {
	return client.GetStatus() == "combat"
}

// ── Dispatch ──

// Dispatch parses a command line and executes the matching command.
func Dispatch(conn *websocket.Conn, entityID string, line string) {
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return
	}

	node, remaining := resolve(tokens)
	if node == nil {
		fmt.Printf("未知命令: %s (输入 help 查看帮助)\n", tokens[0])
		return
	}
	if node.Handler == nil {
		if len(node.Sub) > 0 {
			listChildren(node)
		} else {
			fmt.Printf("未知命令: %s\n", tokens[0])
		}
		return
	}
	node.Handler(conn, entityID, remaining)
}

func resolve(tokens []string) (*CmdDef, []string) {
	for _, c := range cmdTree {
		if matchNames(c.Names, tokens[0]) {
			return resolveSub(c, tokens[1:])
		}
	}
	return nil, nil
}

func resolveSub(parent *CmdDef, tokens []string) (*CmdDef, []string) {
	if len(tokens) == 0 || len(parent.Sub) == 0 {
		return parent, tokens
	}
	for _, c := range parent.Sub {
		if matchNames(c.Names, tokens[0]) {
			if c.Handler != nil {
				return c, tokens[1:]
			}
			return resolveSub(c, tokens[1:])
		}
	}
	return parent, tokens // no sub-match, return parent
}

func matchNames(names []string, token string) bool {
	for _, n := range names {
		if n == token {
			return true
		}
	}
	return false
}

func listChildren(node *CmdDef) {
	fmt.Printf("可用子命令: ")
	for _, c := range node.Sub {
		if c.Context != nil && !c.Context() {
			continue
		}
		fmt.Printf("%s ", c.Names[0])
	}
	fmt.Println()
}

// ── Handlers ──

func wrapOp(actionType string) cmdHandler {
	return func(conn *websocket.Conn, entityID string, args []string) {
		client.SendAction(conn, actionType, nil)
	}
}

func cmdHelp(conn *websocket.Conn, entityID string, args []string) {
	PrintHelp()
}

func cmdClear(conn *websocket.Conn, entityID string, args []string) {
	fmt.Print("\033[H\033[2J")
}

func cmdExit(conn *websocket.Conn, entityID string, args []string) {
	fmt.Println("再见!")
	os.Exit(0)
}

func cmdStatus(conn *websocket.Conn, entityID string, args []string) {
	printStatus()
}

func cmdAttributes(conn *websocket.Conn, entityID string, args []string) {
	entity := client.GetCharacter()
	if entity == nil {
		fmt.Println("角色信息不可用")
		return
	}
	attrs, _ := entity["attributes"].(map[string]interface{})
	if attrs == nil {
		fmt.Println("无属性数据")
		return
	}
	fmt.Println()
	fmt.Printf("  ── 基础 ──\n")
	fmt.Printf("  悟性: %d    根骨: %d    机缘: %d\n",
		getIntDef(attrs, "comprehension"), getIntDef(attrs, "constitution"), getIntDef(attrs, "luck"))
	fmt.Printf("  神识: %d    道心: %d    悟道: %d\n",
		getIntDef(attrs, "divine_sense"), getIntDef(attrs, "dao_heart"), getIntDef(attrs, "enlightenment"))
	fmt.Printf("  寿命: %d / %d\n",
		getIntDef(attrs, "remaining_lifespan"), getIntDef(attrs, "max_lifespan"))
	fmt.Printf("  灵根 purity: %d\n", getIntDef(attrs, "root_purity"))
	if ss, ok := attrs["spirit_stones"].(map[string]interface{}); ok {
		fmt.Printf("  灵石: %d (低) %d (中) %d (高) %d (极)\n",
			getInt64Def(ss, "low_grade"), getInt64Def(ss, "medium_grade"),
			getInt64Def(ss, "high_grade"), getInt64Def(ss, "premium_grade"))
	}
	fmt.Printf("  业力: %d    功德: %d\n",
		getIntDef(attrs, "karma_value"), getIntDef(attrs, "merit"))

	fmt.Printf("\n  ── 战斗 ──\n")
	fmt.Printf("  攻击: %d    防御: %d    速度: %d\n",
		getIntDef(attrs, "attack_power"), getIntDef(attrs, "defense"), getIntDef(attrs, "speed"))
	fmt.Printf("  会心: %.0f%%  会伤: %.0f%%  闪避: %.0f%%  命中: %.0f%%\n",
		getFloatDef(attrs, "crit_rate"), getFloatDef(attrs, "crit_damage"),
		getFloatDef(attrs, "dodge_rate"), getFloatDef(attrs, "hit_rate"))
	fmt.Printf("  穿透: %.0f    减伤: %.0f\n",
		getFloatDef(attrs, "penetration"), getFloatDef(attrs, "damage_reduction"))

	fmt.Printf("\n  ── 炼丹/炼器 ──\n")
	fmt.Printf("  炼丹: %d    炼器: %d    阵法: %d\n",
		getIntDef(attrs, "alchemy_level"), getIntDef(attrs, "artificing_level"), getIntDef(attrs, "formation_level"))
	fmt.Printf("  采药: %d    采矿: %d    制符: %d    训兽: %d    控火: %d\n",
		getIntDef(attrs, "herb_knowledge"), getIntDef(attrs, "mining_skill"),
		getIntDef(attrs, "talisman_skill"), getIntDef(attrs, "beast_taming"),
		getIntDef(attrs, "fire_control"))

	fmt.Printf("\n  ── 社交/状态 ──\n")
	fmt.Printf("  声望: %d    宗门贡献: %d\n",
		getIntDef(attrs, "reputation"), getIntDef(attrs, "sect_contribution"))
	fmt.Printf("  心神: %d\n", getIntDef(attrs, "mental_stability"))
	if pl := getIntDef(attrs, "poison_level"); pl > 0 {
		fmt.Printf("  毒: %d\n", pl)
	}
	if cl := getIntDef(attrs, "curse_level"); cl > 0 {
		fmt.Printf("  咒: %d\n", cl)
	}
	fmt.Println()
}

func cmdSkills(conn *websocket.Conn, entityID string, args []string) {
	spells := client.GetSpells()
	if len(spells) == 0 {
		fmt.Println("未学习任何技能")
		return
	}
	fmt.Println()
	fmt.Printf("  已学技能 (%d):\n", len(spells))
	fmt.Println(strings.Repeat("─", 60))
	for _, s := range spells {
		sp, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		name := getStr(sp, "name")
		if name == "" {
			name = getStr(sp, "spell_id")
		}
		stype := getStr(sp, "spell_type")
		elem := getStr(sp, "element")
		cost := getIntDef(sp, "cost")
		dmg := getIntDef(sp, "base_damage")
		heal := getIntDef(sp, "base_heal")
		cd := getIntDef(sp, "cooldown")
		prof := getIntDef(sp, "proficiency")
		cdRem := getIntDef(sp, "cooldown_remaining")

		line := fmt.Sprintf("  %s", name)
		if stype != "" {
			line += fmt.Sprintf(" [%s", stype)
			if elem != "" {
				line += fmt.Sprintf("/%s", elemDisplay(elem))
			}
			line += "]"
		}
		if dmg > 0 {
			line += fmt.Sprintf(" 伤害:%d", dmg)
		}
		if heal > 0 {
			line += fmt.Sprintf(" 治疗:%d", heal)
		}
		if cost > 0 {
			line += fmt.Sprintf(" 消耗:%d", cost)
		}
		if cd > 0 {
			line += fmt.Sprintf(" 冷却:%ds", cd)
		}
		if prof > 0 {
			line += fmt.Sprintf(" 熟练度:%d", prof)
		}
		if cdRem > 0 {
			line += fmt.Sprintf(" [冷却中:%ds]", cdRem)
		}
		fmt.Println(line)

		desc := getStr(sp, "description")
		if desc != "" {
			fmt.Printf("    └ %s\n", desc)
		}
	}
	fmt.Println()
}

func cmdGather(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 采集 <资源类型> [数量]")
		return
	}
	qty := 1
	if len(args) > 1 {
		if q, err := strconv.Atoi(args[1]); err == nil && q > 0 {
			qty = q
		}
	}
	client.SendAction(conn, "gather", map[string]interface{}{
		"resource_type": args[0],
		"quantity":      qty,
	})
}

func cmdMove(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 移动 <区域ID>")
		return
	}
	client.SendAction(conn, "move", map[string]interface{}{
		"region_id": args[0],
	})
}

func cmdCraft(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 炼制 <配方ID>")
		return
	}
	client.SendAction(conn, "craft", map[string]interface{}{
		"recipe_id": args[0],
	})
}

func cmdCombat(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 目标 <目标ID>")
		return
	}
	client.SendAction(conn, "combat", map[string]interface{}{
		"target_id": args[0],
	})
}

func cmdAutoCombat(conn *websocket.Conn, entityID string, args []string) {
	targetID := entityID
	if len(args) >= 1 {
		targetID = args[0]
	}
	if targetID == entityID {
		fmt.Println("用法: 自动 <目标ID>")
		return
	}

	fmt.Printf("开始自动战斗，目标: %s (按 Ctrl+C 停止)\n", targetID)
	fmt.Println("每3秒自动攻击一次...")

	go func() {
		for {
			client.SendAction(conn, "combat", map[string]interface{}{
				"target_id": targetID,
			})
			time.Sleep(3 * time.Second)

			status := client.GetStatus()
			if status != "combat" {
				fmt.Println("自动战斗结束")
				return
			}
		}
	}()
}

func cmdUseSkill(conn *websocket.Conn, entityID string, args []string) {
	spellID := ""
	targetID := entityID
	if len(args) >= 1 {
		spellID = args[0]
	}
	if len(args) >= 2 {
		targetID = args[1]
	}
	if spellID != "" {
		client.SendAction(conn, "cast_spell", map[string]interface{}{
			"spell_id":  spellID,
			"target_id": targetID,
		})
	} else {
		client.SendAction(conn, "use_skill", nil)
	}
}

func cmdInventoryList(conn *websocket.Conn, entityID string, args []string) {
	items := client.GetItems()
	if len(items) == 0 {
		fmt.Println("背包为空")
		return
	}
	fmt.Println()
	// Group by category
	categories := map[string][]map[string]interface{}{}
	for _, it := range items {
		item, ok := it.(map[string]interface{})
		if !ok {
			continue
		}
		itype := getStr(item, "item_type")
		if itype == "" {
			itype = "其他"
		}
		categories[itype] = append(categories[itype], item)
	}

	for _, catName := range []string{"pill", "material", "weapon", "armor", "talisman", "artifact", "treasure"} {
		catItems, ok := categories[catName]
		if !ok {
			continue
		}
		fmt.Printf("  ── %s (%d) ──\n", itemTypeDisplay(catName), len(catItems))
		for _, item := range catItems {
			name := getStr(item, "name")
			qty := getIntDef(item, "quantity")
			rarity := getIntDef(item, "rarity")
			equipped := false
			if eq, ok := item["equipped"].(bool); ok {
				equipped = eq
			}
			slot := getStr(item, "slot")
			durability := getIntDef(item, "durability")
			bound := false
			if b, ok := item["bound"].(bool); ok {
				bound = b
			}

			line := fmt.Sprintf("    %s", name)
			if qty > 1 {
				line += fmt.Sprintf(" x%d", qty)
			}
			if rarity > 0 {
				line += fmt.Sprintf(" [%s]", rarityDisplay(rarity))
			}
			if bound {
				line += " [已绑定]"
			}
			if equipped {
				line += fmt.Sprintf(" [已装备:%s]", slot)
			}
			if durability > 0 {
				line += fmt.Sprintf(" 耐久:%d", durability)
			}
			fmt.Println(line)

			desc := getStr(item, "description")
			if desc != "" {
				fmt.Printf("      └ %s\n", desc)
			}
		}
		fmt.Println()
	}
	// Remaining categories
	delete(categories, "pill")
	delete(categories, "material")
	delete(categories, "weapon")
	delete(categories, "armor")
	delete(categories, "talisman")
	delete(categories, "artifact")
	delete(categories, "treasure")
	for catName, catItems := range categories {
		fmt.Printf("  ── %s (%d) ──\n", catName, len(catItems))
		for _, item := range catItems {
			name := getStr(item, "name")
			qty := getIntDef(item, "quantity")
			fmt.Printf("    %s x%d\n", name, qty)
		}
		fmt.Println()
	}
	fmt.Printf("  总计: %d 件物品\n", len(items))
}

func cmdUseItem(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 使用 <物品名>")
		return
	}
	itemName := strings.Join(args, " ")
	client.SendAction(conn, "use_item", map[string]interface{}{
		"item_name": itemName,
	})
}

func cmdDropItem(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 丢弃 <物品ID>")
		return
	}
	client.SendAction(conn, "drop_item", map[string]interface{}{
		"item_id": args[0],
	})
}

func cmdChat(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 聊天 <消息>")
		return
	}
	client.SendChat(conn, strings.Join(args, " "), "world")
}

func cmdMessage(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 2 {
		fmt.Println("用法: 私信 <收件人ID> <内容>")
		return
	}
	client.SendAction(conn, "send_message", map[string]interface{}{
		"receiver_id":  args[0],
		"content":      strings.Join(args[1:], " "),
		"message_type": "private",
	})
}

func cmdAddFriend(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 添加 <玩家名>")
		return
	}
	client.SendAction(conn, "add_friend", map[string]interface{}{
		"name": args[0],
	})
}

func cmdRemoveFriend(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 删除 <好友ID>")
		return
	}
	client.SendAction(conn, "remove_friend", map[string]interface{}{
		"friend_id": args[0],
	})
}

func cmdAcceptFriend(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 接受 <请求ID>")
		return
	}
	client.SendAction(conn, "accept_friend", map[string]interface{}{
		"request_id": args[0],
	})
}

func cmdFormSect(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 创建 <宗门名称>")
		return
	}
	client.SendAction(conn, "form_sect", map[string]interface{}{
		"sect_name": strings.Join(args, " "),
	})
}

func cmdJoinSect(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 加入 <宗门ID>")
		return
	}
	client.SendAction(conn, "join_sect", map[string]interface{}{
		"sect_id": args[0],
	})
}

func cmdLeaveSect(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 退出 <宗门ID>")
		return
	}
	client.SendAction(conn, "leave_sect", map[string]interface{}{
		"sect_id": args[0],
	})
}

func cmdTrade(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 3 {
		fmt.Println("用法: 交易 <目标ID> <物品ID> <价格>")
		return
	}
	price, _ := strconv.ParseFloat(args[2], 64)
	client.SendAction(conn, "trade", map[string]interface{}{
		"target_id": args[0],
		"item_id":   args[1],
		"price":     price,
	})
}


// ── New handlers ──

func cmdShowEquipment(conn *websocket.Conn, entityID string, args []string) {
	items := client.GetEquippedItems()
	if len(items) == 0 {
		fmt.Println("未装备任何物品")
		return
	}
	fmt.Println()
	fmt.Printf("  已装备 (%d):\n", len(items))
	fmt.Println(strings.Repeat("─", 60))
	for _, item := range items {
		name := getStr(item, "name")
		slot := getStr(item, "slot")
		rarity := getIntDef(item, "rarity")
		durability := getIntDef(item, "durability")
		line := fmt.Sprintf("  %s [%s]", name, slot)
		if rarity > 0 {
			line += fmt.Sprintf(" [%s]", rarityDisplay(rarity))
		}
		if durability > 0 {
			line += fmt.Sprintf(" 耐久:%d", durability)
		}
		fmt.Println(line)
	}
	fmt.Println()
}

func cmdLearnSpell(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 学习 <法术ID>")
		return
	}
	client.SendAction(conn, "learn_spell", map[string]interface{}{
		"spell_id": args[0],
	})
}

func cmdListFriends(conn *websocket.Conn, entityID string, args []string) {
	// Check if we have cached friends
	friends := client.GetFriends()
	if len(friends) > 0 {
		printFriendList(friends)
		return
	}
	// Request fresh list from server
	client.SendAction(conn, "list_friends", nil)
}

func cmdSectInfo(conn *websocket.Conn, entityID string, args []string) {
	sectID := ""
	if len(args) >= 1 {
		sectID = args[0]
	}
	// Check cache first
	sect := client.GetSect()
	if sectID == "" && sect != nil {
		sectID = getStr(sect, "sect_id")
	}
	if sectID == "" {
		fmt.Println("用法: 信息 <宗门ID>")
		return
	}
	client.SendAction(conn, "sect_info", map[string]interface{}{
		"sect_id": sectID,
	})
}

func printFriendList(friends []interface{}) {
	fmt.Println()
	fmt.Printf("  好友列表 (%d):\n", len(friends))
	fmt.Println(strings.Repeat("─", 60))
	for _, f := range friends {
		fm, ok := f.(map[string]interface{})
		if !ok {
			continue
		}
		name := getStr(fm, "friend_name")
		fid := getStr(fm, "friend_id")
		created := getInt64Def(fm, "created_at")
		line := fmt.Sprintf("  %s", name)
		if name != fid {
			line += fmt.Sprintf(" (%s)", fid)
		}
		if created > 0 {
			t := time.Unix(created, 0)
			line += fmt.Sprintf(" [好友时间:%s]", t.Format("2006-01-02"))
		}
		fmt.Println(line)
	}
	fmt.Println()
}
func cmdEquipItem(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 装备 <物品ID> [装备位]")
		return
	}
	params := map[string]interface{}{
		"item_id": args[0],
	}
	if len(args) >= 2 {
		params["slot"] = args[1]
	}
	client.SendAction(conn, "equip_item", params)
}

func cmdUnequipItem(conn *websocket.Conn, entityID string, args []string) {
	if len(args) < 1 {
		fmt.Println("用法: 卸下 <装备位>")
		return
	}
	client.SendAction(conn, "unequip_item", map[string]interface{}{
		"slot": args[0],
	})
}

// ── Help ──

// PrintHelp prints the grouped command help.
func PrintHelp() {
	fmt.Println()
	for _, c := range cmdTree {
		if c.Context != nil && !c.Context() {
			continue
		}
		if len(c.Desc) == 0 {
			continue // hidden flat alias
		}
		if len(c.Sub) > 0 {
			fmt.Printf("── %s (%s) ──\n", c.Names[0], c.Desc)
			for _, sub := range c.Sub {
				if sub.Context != nil && !sub.Context() {
					continue
				}
				if len(sub.Sub) > 0 {
					for _, s2 := range sub.Sub {
						if s2.Context != nil && !s2.Context() {
							continue
						}
						fmt.Printf("  %s %-18s%s\n", c.Names[0], s2.Names[0], s2.Desc)
					}
				} else {
					u := sub.Usage
					if u == "" {
						u = c.Names[0] + " " + sub.Names[0]
					}
					fmt.Printf("  %-20s%s\n", u, sub.Desc)
				}
			}
		} else if c.Handler != nil {
			u := c.Usage
			if u == "" {
				u = c.Names[0]
			}
			fmt.Printf("  %-20s%s\n", u, c.Desc)
		}
	}
	fmt.Println()
}

// ── Completions ──

// GetCompletions returns context-aware tab-completion candidates for the given prefix.
func GetCompletions(prefix string) []string {
	tokens := strings.Fields(prefix)
	isAfterSpace := strings.HasSuffix(prefix, " ")

	if len(tokens) == 0 || (len(tokens) == 1 && !isAfterSpace) {
		current := ""
		if len(tokens) == 1 {
			current = tokens[0]
		}
		var r []string
		for _, c := range cmdTree {
			if c.Context != nil && !c.Context() {
				continue
			}
			if len(c.Desc) == 0 && len(c.Sub) == 0 {
				continue // hidden flat alias
			}
			for _, n := range c.Names {
				if strings.HasPrefix(n, current) {
					if len(c.Sub) > 0 {
						r = append(r, c.Names[0]+" ") // group → add space
					} else {
						r = append(r, n)
					}
					break
				}
			}
		}
		if len(r) == 0 && current != "" {
			for _, c := range cmdTree {
				if c.Context != nil && !c.Context() {
					continue
				}
				for _, n := range c.Names {
					if strings.HasPrefix(n, current) {
						r = append(r, n)
						break
					}
				}
			}
		}
		return r
	}

	// Completing a subcommand
	group := findGroup(tokens[0])
	if group == nil {
		return nil
	}
	if isAfterSpace && len(tokens) == 1 {
		var r []string
		for _, sub := range group.Sub {
			if sub.Context != nil && !sub.Context() {
				continue
			}
			if len(sub.Sub) > 0 {
				r = append(r, group.Names[0]+" "+sub.Names[0]+" ")
			} else {
				r = append(r, group.Names[0]+" "+sub.Names[0])
			}
		}
		return r
	}
	current := tokens[len(tokens)-1]
	var r []string
	for _, sub := range group.Sub {
		if sub.Context != nil && !sub.Context() {
			continue
		}
		for _, n := range sub.Names {
			if strings.HasPrefix(n, current) {
				if len(sub.Sub) > 0 {
					r = append(r, sub.Names[0]+" ")
				} else {
					r = append(r, n)
				}
				break
			}
		}
	}
	if len(tokens) == 2 && !isAfterSpace {
		for _, sub := range group.Sub {
			if matchNames(sub.Names, tokens[1]) && len(sub.Sub) > 0 {
				return nil
			}
		}
	}
	return r
}

func findGroup(name string) *CmdDef {
	for _, c := range cmdTree {
		if matchNames(c.Names, name) && len(c.Sub) > 0 {
			return c
		}
	}
	return nil
}

// ── Status display ──

func printStatus() {
	entity := client.GetCharacter()
	if entity == nil {
		fmt.Println("角色信息不可用")
		return
	}

	name := getStr(entity, "name")
	realm := getStr(entity, "realm")
	status := getStr(entity, "status")

	attrs, _ := entity["attributes"].(map[string]interface{})

	fmt.Println()
	fmt.Printf("  ╔══════════════════════════════════════╗\n")
	fmt.Printf("  ║  %s\n", centerText(fmt.Sprintf("%s | %s | %s", name, realmDisplay(realm), statusDisplay(status)), 38))
	fmt.Printf("  ╚══════════════════════════════════════╝\n")

	if attrs != nil {
		// HP bars
		if qi, ok := getFloat(attrs, "qi"); ok {
			maxQi, _ := getFloat(attrs, "max_qi")
			fmt.Printf("  灵力: %s %.0f/%.0f\n", hpBar(qi, maxQi, 20), qi, maxQi)
		}
		if sp, ok := getFloat(attrs, "spiritual_power"); ok {
			maxSp, _ := getFloat(attrs, "max_spiritual_power")
			fmt.Printf("  神识: %s %.0f/%.0f\n", hpBar(sp, maxSp, 20), sp, maxSp)
		}
		if prog, ok := getFloat(attrs, "cultivation_progress"); ok {
			fmt.Printf("  修为: %s %.1f%%\n", hpBar(prog, 100, 20), prog)
		}
		if lp, ok := getInt(attrs, "remaining_lifespan"); ok {
			maxLp, _ := getInt(attrs, "max_lifespan")
			fmt.Printf("  寿元: %d / %d\n", lp, maxLp)
		}

		fmt.Println(strings.Repeat("─", 40))

		// Core attributes
		fmt.Printf("  悟性:%-4d 根骨:%-4d 机缘:%-4d\n",
			getIntDef(attrs, "comprehension"), getIntDef(attrs, "constitution"), getIntDef(attrs, "luck"))
		fmt.Printf("  神识:%-4d 道心:%-4d 悟道:%-4d\n",
			getIntDef(attrs, "divine_sense"), getIntDef(attrs, "dao_heart"), getIntDef(attrs, "enlightenment"))
		fmt.Printf("  攻击:%-4d 防御:%-4d 速度:%-4d\n",
			getIntDef(attrs, "attack_power"), getIntDef(attrs, "defense"), getIntDef(attrs, "speed"))

		// Spirit stones
		if ss, ok := attrs["spirit_stones"].(map[string]interface{}); ok {
			lg := getInt64Def(ss, "low_grade")
			mg := getInt64Def(ss, "medium_grade")
			hg := getInt64Def(ss, "high_grade")
			pg := getInt64Def(ss, "premium_grade")
			if lg+mg+hg+pg > 0 {
				fmt.Printf("  灵石: %d低 %d中 %d高 %d极\n", lg, mg, hg, pg)
			}
		}
		fmt.Printf("  业力: %d / 功德: %d\n",
			getIntDef(attrs, "karma_value"), getIntDef(attrs, "merit"))

		// Equipment
		equipped := client.GetEquippedItems()
		if len(equipped) > 0 {
			fmt.Println(strings.Repeat("─", 40))
			fmt.Printf("  装备 (%d):\n", len(equipped))
			for _, item := range equipped {
				iname := getStr(item, "name")
				islot := getStr(item, "slot")
				irarity := getIntDef(item, "rarity")
				line := fmt.Sprintf("    %s [%s]", iname, islot)
				if irarity > 0 {
					line += fmt.Sprintf(" %s", rarityDisplay(irarity))
				}
				fmt.Println(line)
			}
		}
	}

	if pos, ok := entity["position"].(map[string]interface{}); ok {
		rid := getStr(pos, "region_id")
		x, _ := getFloat(pos, "x")
		y, _ := getFloat(pos, "y")
		fmt.Printf("  位置: %s (%.1f, %.1f)\n", rid, x, y)
	}

	fmt.Println()
}

func hpBar(current, max float64, width int) string {
	if max <= 0 {
		return ""
	}
	filled := int((current / max) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"
	return bar
}

func centerText(text string, width int) string {
	// Simple ASCII-length centering
	runes := len([]rune(text))
	if runes >= width || runes == 0 {
		return text
	}
	left := (width - runes) / 2
	result := ""
	for i := 0; i < left; i++ {
		result += " "
	}
	result += text
	return result
}

// ── helpers ──

func getStr(m map[string]interface{}, key string) string {
	if s, ok := m[key].(string); ok {
		return s
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

func getFloatDef(m map[string]interface{}, key string) float64 {
	v, ok := getFloat(m, key)
	if !ok {
		return 0
	}
	return v
}

func getInt(m map[string]interface{}, key string) (int, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	}
	return 0, false
}

func getIntDef(m map[string]interface{}, key string) int {
	v, ok := getInt(m, key)
	if !ok {
		return 0
	}
	return v
}

func getInt64(m map[string]interface{}, key string) (int64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	}
	return 0, false
}

func getInt64Def(m map[string]interface{}, key string) int64 {
	v, ok := getInt64(m, key)
	if !ok {
		return 0
	}
	return v
}

func elemDisplay(e string) string {
	switch e {
	case "fire":
		return "火"
	case "water":
		return "水"
	case "earth":
		return "土"
	case "metal":
		return "金"
	case "wood":
		return "木"
	case "wind":
		return "风"
	case "thunder":
		return "雷"
	case "ice":
		return "冰"
	case "light":
		return "光"
	case "dark":
		return "暗"
	}
	return e
}

func itemTypeDisplay(t string) string {
	switch t {
	case "weapon":
		return "武器"
	case "armor":
		return "防具"
	case "pill":
		return "丹药"
	case "material":
		return "材料"
	case "talisman":
		return "符箓"
	case "artifact":
		return "法宝"
	case "treasure":
		return "宝物"
	}
	return t
}

func rarityDisplay(r int) string {
	switch r {
	case 1:
		return "凡品"
	case 2:
		return "下品"
	case 3:
		return "中品"
	case 4:
		return "上品"
	case 5:
		return "极品"
	}
	return fmt.Sprintf("稀有度%d", r)
}

func realmDisplay(r string) string {
	switch r {
	case "mortal":
		return "凡人"
	case "qi_condensation":
		return "练气期"
	case "foundation":
		return "筑基期"
	case "golden_core":
		return "金丹期"
	case "nascent_soul":
		return "元婴期"
	case "soul_transformation":
		return "化神期"
	case "void_refinement":
		return "炼虚期"
	case "integration":
		return "合体期"
	case "mahayana":
		return "大乘期"
	case "tribulation":
		return "渡劫期"
	}
	return r
}

func statusDisplay(s string) string {
	switch s {
	case "normal":
		return "正常"
	case "cultivating":
		return "修炼中"
	case "combat":
		return "战斗中"
	case "resting":
		return "休息中"
	case "dead":
		return "已死亡"
	case "exploring":
		return "探索中"
	case "crafting":
		return "炼制中"
	case "meditating":
		return "打坐中"
	default:
		if s == "" {
			return "正常"
		}
		return s
	}
}
