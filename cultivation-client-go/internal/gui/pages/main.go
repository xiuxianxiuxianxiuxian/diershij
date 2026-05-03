package pages

import (
	"cultivation-client/internal/gui/components"
	"cultivation-client/internal/gui/pages/tabs"
	"gioui.org/layout"
	"gioui.org/unit"
)

type MainPage struct {
	sidebar             *components.Sidebar
	statusBar           *components.StatusBar
	notificationManager *components.NotificationManager
	tabs                map[string]Tab
	currentTabID        string
	onTabChange         func(string)
}

type Tab interface {
	Layout(gtx layout.Context) layout.Dimensions
}

func NewMainPage() *MainPage {
	mp := &MainPage{
		sidebar:             components.NewSidebar(),
		statusBar:           components.NewStatusBar(),
		notificationManager: components.NewNotificationManager(),
		tabs: make(map[string]Tab),
		currentTabID:        "character",
	}

	mp.tabs["character"] = tabs.NewCharacterTab()
	mp.tabs["combat"] = tabs.NewCombatTab()
	mp.tabs["world"] = tabs.NewWorldTab()
	mp.tabs["social"] = tabs.NewSocialTab()
	mp.tabs["settings"] = tabs.NewSettingsTab()

	return mp
}

func (mp *MainPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(40))
			return mp.statusBar.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))
					return mp.sidebar.Layout(gtx, mp.currentTabID, mp.handleTabChange)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return mp.layoutContent(gtx)
				}),
			)
		}),
	)
}

func (mp *MainPage) layoutContent(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return mp.notificationManager.Layout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			tab, ok := mp.tabs[mp.currentTabID]
			if !ok {
				return layout.Dimensions{}
			}
			return tab.Layout(gtx)
		}),
	)
}

func (mp *MainPage) handleTabChange(tabID string) {
	mp.currentTabID = tabID
	if mp.onTabChange != nil {
		mp.onTabChange(tabID)
	}
}

func (mp *MainPage) GetNotificationManager() *components.NotificationManager {
	return mp.notificationManager
}

func (mp *MainPage) SetOnTabChange(callback func(string)) {
	mp.onTabChange = callback
}

func (mp *MainPage) GetCurrentTabID() string {
	return mp.currentTabID
}

func (mp *MainPage) SetCurrentTabID(tabID string) {
	mp.currentTabID = tabID
}