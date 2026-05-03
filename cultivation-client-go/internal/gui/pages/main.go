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
	tabs                map[string