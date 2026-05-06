package commands

import (
	"strings"
)

// ActionType constants for game server operations.
const (
	ActionCultivate     = "cultivate"
	ActionBreakthrough  = "breakthrough"
	ActionCombat        = "combat"
	ActionExplore       = "explore"
	ActionGather        = "gather"
	ActionCraft         = "craft"
	ActionCreateMethod  = "create_method"
	ActionMeditate      = "meditate"
	ActionSleep         = "sleep"
	ActionMove          = "move"
	ActionFlee          = "flee"
	ActionUseSkill      = "use_skill"
	ActionCastSpell     = "cast_spell"
	ActionTrade         = "trade"
	ActionFormSect      = "form_sect"
	ActionJoinSect      = "join_sect"
	ActionLeaveSect     = "leave_sect"
	ActionSendMessage   = "send_message"
	ActionAddFriend     = "add_friend"
	ActionRemoveFriend  = "remove_friend"
	ActionAcceptFriend  = "accept_friend"
	ActionLearnSpell    = "learn_spell"
	ActionEquipItem     = "equip_item"
	ActionUnequipItem   = "unequip_item"
	ActionUseItem       = "use_item"
	ActionDropItem      = "drop_item"
	ActionMethodLearn   = "learn_method"
	ActionMethodSetMain = "set_main_method"
	ActionShopList      = "shop_list"
	ActionShopItems     = "shop_items"
	ActionShopBuy       = "buy"
	ActionShopSell      = "sell"
	ActionAuctionList   = "auction_list"
	ActionAuctionBuy    = "auction_buy"
	ActionAuctionCancel = "auction_cancel"
	ActionAuctionView   = "auction_view"
	ActionMailList      = "mail_list"
	ActionMailRead      = "mail_read"
	ActionMailClaim     = "mail_claim"
	ActionMailDelete    = "mail_delete"
	ActionLeaderboard   = "leaderboard"
	ActionNearby        = "nearby_players"
	ActionSectInfo      = "sect_info"
	ActionAutoCombat    = "auto_combat"
	ActionWorldEvents   = "world_events"
	ActionListFriends   = "list_friends"
	ActionListMethods   = "list_methods"
)

// Command represents a parsed user command.
type Command struct {
	Action     string                 // "builtin", "op", "chat", "quit"
	Name       string                 // command name for builtins
	ActionType string                 // operation type for game actions
	Params     map[string]interface{} // parameters for game actions
	Content    string                 // chat content
	Raw        string                 // original line
}

// Parse parses a command line and returns a structured Command, or nil if unknown.
func Parse(line string) *Command {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return nil
	}

	cmd := &Command{Raw: line}

	switch tokens[0] {
	// ── Built-in commands ──
	case "help":
		cmd.Action = "builtin"
		cmd.Name = "help"
		return cmd

	case "clear", "cls":
		cmd.Action = "builtin"
		cmd.Name = "clear"
		return cmd

	case "exit", "quit":
		cmd.Action = "quit"
		return cmd

	case "map":
		cmd.Action = "builtin"
		cmd.Name = "map"
		return cmd

	// ── Character ──
	case "status", "st":
		cmd.Action = "builtin"
		cmd.Name = "status"
		return cmd

	case "attributes", "attrs":
		cmd.Action = "builtin"
		cmd.Name = "attributes"
		return cmd

	case "skills":
		cmd.Action = "builtin"
		cmd.Name = "skills"
		return cmd

	// ── Combat ──
	case "attack", "atk":
		if len(tokens) < 2 {
			return &Command{Action: "op", ActionType: ActionCombat, Params: map[string]interface{}{}}
		}
		return &Command{
			Action:     "op",
			ActionType: ActionCombat,
			Params:     map[string]interface{}{"target_id": tokens[1]},
		}

	case "auto", "auto_combat":
		return &Command{Action: "op", ActionType: ActionAutoCombat, Params: map[string]interface{}{}}

	case "flee":
		return &Command{Action: "op", ActionType: ActionFlee, Params: map[string]interface{}{}}

	case "skill", "use_skill":
		if len(tokens) < 2 {
			return nil
		}
		params := map[string]interface{}{"skill_id": tokens[1]}
		if len(tokens) >= 3 {
			params["target_id"] = tokens[2]
		}
		return &Command{Action: "op", ActionType: ActionUseSkill, Params: params}

	case "cast_spell", "spell":
		if len(tokens) < 2 {
			return nil
		}
		params := map[string]interface{}{"spell_id": tokens[1]}
		if len(tokens) >= 3 {
			params["target_id"] = tokens[2]
		}
		return &Command{Action: "op", ActionType: ActionCastSpell, Params: params}

	// ── Cultivation ──
	case "cultivate", "cult":
		return &Command{Action: "op", ActionType: ActionCultivate, Params: map[string]interface{}{}}

	case "meditate", "med":
		return &Command{Action: "op", ActionType: ActionMeditate, Params: map[string]interface{}{}}

	case "sleep":
		return &Command{Action: "op", ActionType: ActionSleep, Params: map[string]interface{}{}}

	case "breakthrough", "bt":
		return &Command{Action: "op", ActionType: ActionBreakthrough, Params: map[string]interface{}{}}

	case "explore", "exp":
		return &Command{Action: "op", ActionType: ActionExplore, Params: map[string]interface{}{}}

	case "gather":
		params := map[string]interface{}{}
		if len(tokens) >= 2 {
			params["resource_type"] = tokens[1]
		}
		if len(tokens) >= 3 {
			params["quantity"] = tokens[2]
		}
		return &Command{Action: "op", ActionType: ActionGather, Params: params}

	case "move":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionMove,
			Params:     map[string]interface{}{"region_id": tokens[1]},
		}

	case "craft":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionCraft,
			Params:     map[string]interface{}{"recipe_id": tokens[1]},
		}

	case "create_method", "cm":
		return &Command{Action: "op", ActionType: ActionCreateMethod, Params: map[string]interface{}{}}

	case "learn_spell", "learn":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionLearnSpell,
			Params:     map[string]interface{}{"spell_id": tokens[1]},
		}

	// ── Chat ──
	case "chat":
		if len(tokens) < 2 {
			return nil
		}
		cmd.Action = "chat"
		cmd.Content = strings.Join(tokens[1:], " ")
		return cmd

	case "msg", "message", "send_message":
		if len(tokens) < 3 {
			return nil
		}
		cmd.Action = "op"
		cmd.ActionType = ActionSendMessage
		cmd.Params = map[string]interface{}{
			"target_id": tokens[1],
			"content":   strings.Join(tokens[2:], " "),
		}
		return cmd

	// ── Social ──
	case "friend", "add_friend":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionAddFriend,
			Params:     map[string]interface{}{"target_name": tokens[1]},
		}

	case "remove_friend", "unfriend":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionRemoveFriend,
			Params:     map[string]interface{}{"friend_id": tokens[1]},
		}

	case "accept_friend", "accept":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionAcceptFriend,
			Params:     map[string]interface{}{"request_id": tokens[1]},
		}

	case "list_friends":
		return &Command{Action: "op", ActionType: ActionListFriends, Params: map[string]interface{}{}}

	// ── Sect ──
	case "form_sect", "create_sect":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionFormSect,
			Params:     map[string]interface{}{"name": strings.Join(tokens[1:], " ")},
		}

	case "join_sect":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionJoinSect,
			Params:     map[string]interface{}{"sect_id": tokens[1]},
		}

	case "leave_sect":
		params := map[string]interface{}{}
		if len(tokens) >= 2 {
			params["sect_id"] = tokens[1]
		}
		return &Command{Action: "op", ActionType: ActionLeaveSect, Params: params}

	case "sect_info":
		params := map[string]interface{}{}
		if len(tokens) >= 2 {
			params["sect_id"] = tokens[1]
		}
		return &Command{Action: "op", ActionType: ActionSectInfo, Params: params}

	// ── Trade ──
	case "trade":
		if len(tokens) < 4 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionTrade,
			Params: map[string]interface{}{
				"target_id": tokens[1],
				"item_id":   tokens[2],
				"price":     tokens[3],
			},
		}

	// ── Equipment ──
	case "equip":
		if len(tokens) < 2 {
			return &Command{Action: "builtin", Name: "equip"}
		}
		params := map[string]interface{}{"item_id": tokens[1]}
		if len(tokens) >= 3 {
			params["slot"] = tokens[2]
		}
		return &Command{Action: "op", ActionType: ActionEquipItem, Params: params}

	case "unequip":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionUnequipItem,
			Params:     map[string]interface{}{"slot": tokens[1]},
		}

	// ── Inventory ──
	case "inventory", "bag":
		cmd.Action = "builtin"
		cmd.Name = "inventory"
		return cmd

	case "use":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionUseItem,
			Params:     map[string]interface{}{"item_name": strings.Join(tokens[1:], " ")},
		}

	case "drop":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionDropItem,
			Params:     map[string]interface{}{"item_id": tokens[1]},
		}

	// ── Methods ──
	case "list_methods":
		return &Command{Action: "op", ActionType: ActionListMethods, Params: map[string]interface{}{}}

	case "learn_method":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionMethodLearn,
			Params:     map[string]interface{}{"method_id": tokens[1]},
		}

	case "set_main_method":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionMethodSetMain,
			Params:     map[string]interface{}{"method_id": tokens[1]},
		}

	// ── Shop ──
	case "shop_list":
		return &Command{Action: "op", ActionType: ActionShopList, Params: map[string]interface{}{}}

	case "shop_items":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionShopItems,
			Params:     map[string]interface{}{"shop_id": tokens[1]},
		}

	case "buy":
		if len(tokens) < 3 {
			return nil
		}
		params := map[string]interface{}{
			"shop_id":   tokens[1],
			"item_name": tokens[2],
		}
		if len(tokens) >= 4 {
			params["quantity"] = tokens[3]
		}
		return &Command{Action: "op", ActionType: ActionShopBuy, Params: params}

	case "sell":
		if len(tokens) < 3 {
			return nil
		}
		params := map[string]interface{}{
			"shop_id":   tokens[1],
			"item_name": tokens[2],
		}
		if len(tokens) >= 4 {
			params["quantity"] = tokens[3]
		}
		return &Command{Action: "op", ActionType: ActionShopSell, Params: params}

	// ── Auction ──
	case "auction_list":
		return &Command{Action: "op", ActionType: ActionAuctionView, Params: map[string]interface{}{}}

	case "auction_buy":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionAuctionBuy,
			Params:     map[string]interface{}{"auction_id": tokens[1]},
		}

	case "auction_create":
		if len(tokens) < 3 {
			return nil
		}
		params := map[string]interface{}{
			"item_name": tokens[1],
			"price":     tokens[2],
		}
		if len(tokens) >= 4 {
			params["quantity"] = tokens[3]
		}
		return &Command{Action: "op", ActionType: ActionAuctionList, Params: params}

	case "auction_cancel":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionAuctionCancel,
			Params:     map[string]interface{}{"auction_id": tokens[1]},
		}

	case "auction_view":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action:     "op",
			ActionType: ActionAuctionView,
			Params:     map[string]interface{}{"auction_id": tokens[1]},
		}

	// ── Leaderboard ──
	case "leaderboard", "rank":
		boardType := "cultivation"
		if len(tokens) >= 2 {
			boardType = tokens[1]
		}
		return &Command{
			Action:     "op",
			ActionType: ActionLeaderboard,
			Params:     map[string]interface{}{"board_type": boardType},
		}

	// ── Mail ──
	case "mail":
		if len(tokens) < 2 {
			return nil
		}
		switch tokens[1] {
		case "list", "ls":
			return &Command{Action: "op", ActionType: ActionMailList, Params: map[string]interface{}{}}
		case "read":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMailRead,
				Params: map[string]interface{}{"mail_id": tokens[2]},
			}
		case "claim":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMailClaim,
				Params: map[string]interface{}{"mail_id": tokens[2]},
			}
		case "delete", "del":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMailDelete,
				Params: map[string]interface{}{"mail_id": tokens[2]},
			}
		}
		return nil

	// ── Nearby ──
	case "nearby", "near":
		return &Command{Action: "op", ActionType: ActionNearby, Params: map[string]interface{}{}}

	// ── World Events ──
	case "world_events", "events":
		return &Command{Action: "op", ActionType: ActionWorldEvents, Params: map[string]interface{}{}}

	// ══════════════════════════════════════════════
	// Chinese command aliases
	// ══════════════════════════════════════════════

	// ── Built-in ──
	case "帮助":
		cmd.Action = "builtin"
		cmd.Name = "help"
		return cmd
	case "清屏":
		cmd.Action = "builtin"
		cmd.Name = "clear"
		return cmd
	case "退出":
		cmd.Action = "quit"
		return cmd
	case "状态":
		cmd.Action = "builtin"
		cmd.Name = "status"
		return cmd
	case "属性":
		cmd.Action = "builtin"
		cmd.Name = "attributes"
		return cmd
	case "技能":
		cmd.Action = "builtin"
		cmd.Name = "skills"
		return cmd
	case "背包":
		cmd.Action = "builtin"
		cmd.Name = "inventory"
		return cmd
	case "地图":
		cmd.Action = "builtin"
		cmd.Name = "map"
		return cmd

	// ── Cultivation ──
	case "修炼":
		return &Command{Action: "op", ActionType: ActionCultivate, Params: map[string]interface{}{}}
	case "打坐":
		return &Command{Action: "op", ActionType: ActionMeditate, Params: map[string]interface{}{}}
	case "休息":
		return &Command{Action: "op", ActionType: ActionSleep, Params: map[string]interface{}{}}
	case "突破":
		return &Command{Action: "op", ActionType: ActionBreakthrough, Params: map[string]interface{}{}}

	// ── Exploration ──
	case "探索":
		return &Command{Action: "op", ActionType: ActionExplore, Params: map[string]interface{}{}}
	case "采集":
		p := map[string]interface{}{}
		if len(tokens) >= 2 {
			p["resource_type"] = tokens[1]
		}
		if len(tokens) >= 3 {
			p["quantity"] = tokens[2]
		}
		return &Command{Action: "op", ActionType: ActionGather, Params: p}
	case "移动":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action: "op", ActionType: ActionMove,
			Params: map[string]interface{}{"region_id": tokens[1]},
		}
	case "炼制":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action: "op", ActionType: ActionCraft,
			Params: map[string]interface{}{"recipe_id": tokens[1]},
		}
	case "自创":
		return &Command{Action: "op", ActionType: ActionCreateMethod, Params: map[string]interface{}{}}

	// ── Combat ──
	case "战斗":
		return &Command{Action: "op", ActionType: ActionCombat, Params: map[string]interface{}{}}
	case "攻击":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action: "op", ActionType: ActionCombat,
			Params: map[string]interface{}{"target_id": tokens[1]},
		}
	case "自动", "自动战斗":
		return &Command{Action: "op", ActionType: ActionAutoCombat, Params: map[string]interface{}{}}
	case "逃跑":
		return &Command{Action: "op", ActionType: ActionFlee, Params: map[string]interface{}{}}
	case "法术", "施展":
		if len(tokens) < 2 {
			return nil
		}
		params := map[string]interface{}{"spell_id": tokens[1]}
		if len(tokens) >= 3 {
			params["target_id"] = tokens[2]
		}
		return &Command{Action: "op", ActionType: ActionCastSpell, Params: params}

	// ── Inventory ──
	case "使用":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action: "op", ActionType: ActionUseItem,
			Params: map[string]interface{}{"item_name": strings.Join(tokens[1:], " ")},
		}
	case "丢弃":
		if len(tokens) < 2 {
			return nil
		}
		return &Command{
			Action: "op", ActionType: ActionDropItem,
			Params: map[string]interface{}{"item_id": tokens[1]},
		}

	// ── Equipment ──
	case "装备":
		if len(tokens) >= 2 && (tokens[1] == "列表" || tokens[1] == "list" || tokens[1] == "ls") {
			cmd.Action = "builtin"
			cmd.Name = "equip"
			return cmd
		}
		if len(tokens) >= 3 && tokens[1] == "装备" {
			params := map[string]interface{}{"item_id": tokens[2]}
			if len(tokens) >= 4 {
				params["slot"] = tokens[3]
			}
			return &Command{Action: "op", ActionType: ActionEquipItem, Params: params}
		}
		if len(tokens) >= 3 && tokens[1] == "卸下" {
			return &Command{
				Action: "op", ActionType: ActionUnequipItem,
				Params: map[string]interface{}{"slot": tokens[2]},
			}
		}
		cmd.Action = "builtin"
		cmd.Name = "equip"
		return cmd

	// ── Chat ──
	case "聊天":
		if len(tokens) < 2 {
			return nil
		}
		cmd.Action = "chat"
		cmd.Content = strings.Join(tokens[1:], " ")
		return cmd
	case "私信":
		if len(tokens) < 3 {
			return nil
		}
		cmd.Action = "op"
		cmd.ActionType = ActionSendMessage
		cmd.Params = map[string]interface{}{
			"target_id": tokens[1],
			"content":   strings.Join(tokens[2:], " "),
		}
		return cmd

	// ── Friends ──
	case "好友":
		if len(tokens) < 2 {
			return &Command{Action: "op", ActionType: ActionListFriends, Params: map[string]interface{}{}}
		}
		switch tokens[1] {
		case "添加", "add":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionAddFriend,
				Params: map[string]interface{}{"target_name": tokens[2]},
			}
		case "删除", "remove", "del":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionRemoveFriend,
				Params: map[string]interface{}{"friend_id": tokens[2]},
			}
		case "接受", "accept":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionAcceptFriend,
				Params: map[string]interface{}{"request_id": tokens[2]},
			}
		case "列表", "list", "ls":
			return &Command{Action: "op", ActionType: ActionListFriends, Params: map[string]interface{}{}}
		}
		return nil

	// ── Sect ──
	case "宗门":
		if len(tokens) < 2 {
			return &Command{Action: "op", ActionType: ActionSectInfo, Params: map[string]interface{}{}}
		}
		switch tokens[1] {
		case "创建", "create", "form":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionFormSect,
				Params: map[string]interface{}{"name": strings.Join(tokens[2:], " ")},
			}
		case "加入", "join":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionJoinSect,
				Params: map[string]interface{}{"sect_id": tokens[2]},
			}
		case "离开", "leave":
			params := map[string]interface{}{}
			if len(tokens) >= 3 {
				params["sect_id"] = tokens[2]
			}
			return &Command{Action: "op", ActionType: ActionLeaveSect, Params: params}
		case "信息", "info":
			params := map[string]interface{}{}
			if len(tokens) >= 3 {
				params["sect_id"] = tokens[2]
			}
			return &Command{Action: "op", ActionType: ActionSectInfo, Params: params}
		}
		return nil

	// ── Shop ──
	case "商店":
		if len(tokens) < 2 {
			return &Command{Action: "op", ActionType: ActionShopList, Params: map[string]interface{}{}}
		}
		switch tokens[1] {
		case "列表", "list", "ls":
			return &Command{Action: "op", ActionType: ActionShopList, Params: map[string]interface{}{}}
		case "物品", "items":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionShopItems,
				Params: map[string]interface{}{"shop_id": tokens[2]},
			}
		case "购买", "buy":
			if len(tokens) < 4 {
				return nil
			}
			params := map[string]interface{}{
				"shop_id":   tokens[2],
				"item_name": tokens[3],
			}
			if len(tokens) >= 5 {
				params["quantity"] = tokens[4]
			}
			return &Command{Action: "op", ActionType: ActionShopBuy, Params: params}
		case "出售", "sell":
			if len(tokens) < 4 {
				return nil
			}
			params := map[string]interface{}{
				"shop_id":   tokens[2],
				"item_name": tokens[3],
			}
			if len(tokens) >= 5 {
				params["quantity"] = tokens[4]
			}
			return &Command{Action: "op", ActionType: ActionShopSell, Params: params}
		}
		return nil

	// ── Auction ──
	case "拍卖":
		if len(tokens) < 2 {
			return &Command{Action: "op", ActionType: ActionAuctionView, Params: map[string]interface{}{}}
		}
		switch tokens[1] {
		case "列表", "list", "ls":
			return &Command{Action: "op", ActionType: ActionAuctionView, Params: map[string]interface{}{}}
		case "购买", "buy":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionAuctionBuy,
				Params: map[string]interface{}{"auction_id": tokens[2]},
			}
		case "上架", "create":
			if len(tokens) < 4 {
				return nil
			}
			params := map[string]interface{}{
				"item_name": tokens[2],
				"price":     tokens[3],
			}
			if len(tokens) >= 5 {
				params["quantity"] = tokens[4]
			}
			return &Command{Action: "op", ActionType: ActionAuctionList, Params: params}
		case "取消", "cancel":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionAuctionCancel,
				Params: map[string]interface{}{"auction_id": tokens[2]},
			}
		case "查看", "view":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionAuctionView,
				Params: map[string]interface{}{"auction_id": tokens[2]},
			}
		}
		return nil
	// ── Methods ──
	case "功法":
		if len(tokens) < 2 {
			return &Command{Action: "op", ActionType: ActionListMethods, Params: map[string]interface{}{}}
		}
		switch tokens[1] {
		case "列表", "list", "ls":
			return &Command{Action: "op", ActionType: ActionListMethods, Params: map[string]interface{}{}}
		case "学习", "learn":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMethodLearn,
				Params: map[string]interface{}{"method_id": tokens[2]},
			}
		case "主修", "main", "set":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMethodSetMain,
				Params: map[string]interface{}{"method_id": tokens[2]},
			}
		}
		return nil

	// ── Trade ──
	case "交易":
		if len(tokens) < 4 {
			return nil
		}
		return &Command{
			Action: "op", ActionType: ActionTrade,
			Params: map[string]interface{}{
				"target_id": tokens[1],
				"item_id":   tokens[2],
				"price":     tokens[3],
			},
		}

	// ── Mail ──
	case "邮件":
		if len(tokens) < 2 {
			return nil
		}
		switch tokens[1] {
		case "列表", "list", "ls":
			return &Command{Action: "op", ActionType: ActionMailList, Params: map[string]interface{}{}}
		case "读取", "read":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMailRead,
				Params: map[string]interface{}{"mail_id": tokens[2]},
			}
		case "领取", "claim":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMailClaim,
				Params: map[string]interface{}{"mail_id": tokens[2]},
			}
		case "删除", "delete", "del":
			if len(tokens) < 3 {
				return nil
			}
			return &Command{
				Action: "op", ActionType: ActionMailDelete,
				Params: map[string]interface{}{"mail_id": tokens[2]},
			}
		}
		return nil

	// ── Leaderboard ──
	case "排行榜":
		boardType := "cultivation"
		if len(tokens) >= 2 {
			boardType = tokens[1]
		}
		return &Command{
			Action: "op", ActionType: ActionLeaderboard,
			Params: map[string]interface{}{"board_type": boardType},
		}

	// ── Nearby ──
	case "附近":
		return &Command{Action: "op", ActionType: ActionNearby, Params: map[string]interface{}{}}

	// ── World Events ──
	case "世界事件":
		return &Command{Action: "op", ActionType: ActionWorldEvents, Params: map[string]interface{}{}}
	}

	return nil
}

// Completions returns possible completions for a given prefix.
func Completions(prefix string) []string {
	all := []string{
		"help", "clear", "exit", "map",
		"status", "attributes", "skills",
		"cultivate", "meditate", "sleep", "breakthrough", "explore",
		"gather", "move", "craft", "create_method",
		"attack", "auto", "flee", "skill", "cast_spell",
		"inventory", "use", "drop",
		"chat", "message",
		"friend", "add_friend", "remove_friend", "accept_friend", "list_friends",
		"form_sect", "join_sect", "leave_sect", "sect_info",
		"trade",
		"shop_list", "shop_items", "buy", "sell",
		"auction_list", "auction_buy", "auction_create", "auction_cancel", "auction_view",
		"equip", "unequip",
		"list_methods", "learn_method", "set_main_method",
		"leaderboard", "mail", "nearby",
		"world_events", "events",
	}
	if prefix == "" {
		return all
	}
	var result []string
	for _, c := range all {
		if strings.HasPrefix(c, prefix) {
			result = append(result, c)
		}
	}
	return result
}
