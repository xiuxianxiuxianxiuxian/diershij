package app

import (
	"encoding/json"
	"log"

	"cultivation-client/internal/gui/pages"
	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/network"
	"cultivation-client/internal/store"
	"cultivation-client/internal/types"
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

type CultivationApp struct {
	window       *app.Window
	currentView  string
	loginPage    *pages.LoginPage
	registerPage *pages.RegisterPage
	mainPage     *pages.MainPage
	theme        *theme.Theme
}

func New() *CultivationApp {
	th := theme.DefaultTheme

	app := &CultivationApp{
		currentView:  "login",
		loginPage:    pages.NewLoginPage(),
		registerPage: pages.NewRegisterPage(),
		mainPage:     pages.NewMainPage(),
		theme:        &th,
	}

	// 设置标签页切换回调，触发窗口重绘
	app.mainPage.SetOnTabChange(func(tabID string) {
		if app.window != nil {
			app.window.Invalidate()
		}
	})

	// 设置登录成功回调
	app.loginPage.SetOnLogin(func() {
		app.navigateTo("main")
		// 显示登录成功通知
		app.mainPage.GetNotificationManager().AddSuccess("登录成功，欢迎回来！")
	})

	// 设置切换到注册页面回调
	app.loginPage.SetOnRegister(func() {
		app.navigateTo("register")
	})

	// 设置注册成功回调
	app.registerPage.SetOnRegister(func() {
		app.navigateTo("main")
		// 显示注册成功通知
		app.mainPage.GetNotificationManager().AddSuccess("注册成功，欢迎加入修仙世界！")
	})

	// 设置切换到登录页面回调
	app.registerPage.SetOnLogin(func() {
		app.navigateTo("login")
	})

	// 注册所有 WebSocket 处理器（数据更新 + 通知）
	app.registerHandlers()

	return app
}

// registerHandlers 注册所有 WebSocket 消息处理器（数据更新 + 通知）
func (a *CultivationApp) registerHandlers() {
	ws := network.GetWebSocketClient()

	// combat_update：战斗状态更新
	ws.RegisterHandler("combat_update", func(payload []byte) {
		var update types.CombatState
		if err := json.Unmarshal(payload, &update); err != nil {
			return
		}
		store.GetGameStore().SetCombat(&update)
		if update.InCombat && update.CurrentEnemy != nil {
			a.mainPage.GetNotificationManager().AddWarning("遭遇 " + update.CurrentEnemy.Name + "，战斗开始！")
		}
	})

	// world_update：世界状态更新
	ws.RegisterHandler("world_update", func(payload []byte) {
		var update types.WorldState
		if err := json.Unmarshal(payload, &update); err != nil {
			return
		}
		store.GetGameStore().SetWorld(&update)
	})

	// social_update：社交信息更新
	ws.RegisterHandler("social_update", func(payload []byte) {
		var update types.SocialInfo
		if err := json.Unmarshal(payload, &update); err != nil {
			return
		}
		store.GetGameStore().SetSocial(&update)
	})

	// new_message：新消息
	ws.RegisterHandler("new_message", func(payload []byte) {
		var msg types.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			return
		}
		store.GetGameStore().AddMessage(msg)
		a.mainPage.GetNotificationManager().AddInfo("收到新消息来自 " + msg.SenderName)
	})

	// op_result：操作结果
	ws.RegisterHandler("op_result", func(payload []byte) {
		var result types.OperationResult
		if err := json.Unmarshal(payload, &result); err != nil {
			return
		}
		store.GetGameStore().SetLastOperationResult(&result)
		if result.Success {
			a.mainPage.GetNotificationManager().AddSuccess(result.Message)
		} else {
			a.mainPage.GetNotificationManager().AddError(result.Message)
		}
	})

	// state_sync：状态同步（初始全量同步）
	ws.RegisterHandler("state_sync", func(payload []byte) {
		var syncData struct {
			RawEntity map[string]interface{} `json:"entity"`
		}
		if err := json.Unmarshal(payload, &syncData); err != nil || syncData.RawEntity == nil {
			return
		}
		store.GetGameStore().SetCharacterFromServerMap(syncData.RawEntity)
	})

	// entity_update：实体增量更新
	ws.RegisterHandler("entity_update", func(payload []byte) {
		var update struct {
			RawEntity map[string]interface{} `json:"entity"`
		}
		if err := json.Unmarshal(payload, &update); err != nil || update.RawEntity == nil {
			return
		}
		store.GetGameStore().SetCharacterFromServerMap(update.RawEntity)
	})

	// chat：聊天消息
	ws.RegisterHandler("chat", func(payload []byte) {
		var msg struct {
			SenderID   string `json:"sender_id"`
			SenderName string `json:"sender_name"`
			Content    string `json:"content"`
			Channel    string `json:"channel"`
		}
		if err := json.Unmarshal(payload, &msg); err != nil {
			return
		}
		store.GetGameStore().AddMessage(types.Message{
			ID:         msg.SenderID,
			SenderID:   msg.SenderID,
			SenderName: msg.SenderName,
			Content:    msg.Content,
		})
	})

	// disconnect：连接断开
	ws.RegisterHandler("disconnect", func(payload []byte) {
		a.mainPage.GetNotificationManager().AddWarning("与服务器的连接已断开，正在尝试重连...")
	})

	// connected：连接恢复
	ws.RegisterHandler("connected", func(payload []byte) {
		a.mainPage.GetNotificationManager().AddSuccess("已重新连接到服务器")
	})

	// world_event：世界事件
	ws.RegisterHandler("world_event", func(payload []byte) {
		var event types.WorldEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			return
		}
		a.mainPage.GetNotificationManager().AddInfo("世界事件: " + event.Description)
	})

	// announcement：公告
	ws.RegisterHandler("announcement", func(payload []byte) {
		var announcement types.Announcement
		if err := json.Unmarshal(payload, &announcement); err != nil {
			return
		}
		if announcement.Priority >= 2 {
			a.mainPage.GetNotificationManager().AddWarning("公告: " + announcement.Content)
		} else {
			a.mainPage.GetNotificationManager().AddInfo("公告: " + announcement.Content)
		}
	})

	// friend_request：好友请求
	ws.RegisterHandler("friend_request", func(payload []byte) {
		var request types.FriendRequest
		if err := json.Unmarshal(payload, &request); err != nil {
			return
		}
		a.mainPage.GetNotificationManager().AddInfo(request.FromName + " 向你发送了好友请求")
	})

	// system_message：系统消息
	ws.RegisterHandler("system_message", func(payload []byte) {
		var data map[string]interface{}
		if err := json.Unmarshal(payload, &data); err != nil {
			return
		}
		if msg, ok := data["message"].(string); ok {
			a.mainPage.GetNotificationManager().AddInfo(msg)
		}
	})
}

func (a *CultivationApp) Run() error {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("修仙世界"),
			app.Size(unit.Dp(1280), unit.Dp(800)),
		)
		a.window = w
		err := a.run(w)
		if err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
	return nil
}

func (a *CultivationApp) run(w *app.Window) error {
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&op.Ops{}, e)

			switch a.currentView {
			case "login":
				a.loginPage.Layout(gtx)
			case "register":
				a.registerPage.Layout(gtx)
			case "main":
				a.mainPage.Layout(gtx)
			}

			e.Frame(gtx.Ops)
		}
	}
}

func (a *CultivationApp) navigateTo(view string) {
	a.currentView = view
	if a.window != nil {
		a.window.Invalidate()
	}
}

func (a *CultivationApp) handleLogout() {
	network.GetWebSocketClient().Disconnect()
	store.GetAuthStore().Logout()
	store.GetGameStore().Clear()
	a.loginPage.Clear()
	a.navigateTo("login")
}

func (a *CultivationApp) SetTheme(t *theme.Theme) {
	a.theme = t
}
