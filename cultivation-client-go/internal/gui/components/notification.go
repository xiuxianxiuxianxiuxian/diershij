package components

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"cultivation-client/internal/gui/theme"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// NotificationType 通知类型
type NotificationType int

const (
	NotificationSuccess NotificationType = iota
	NotificationError
	NotificationWarning
	NotificationInfo
)

// Notification 单个通知
type Notification struct {
	ID        string
	Message   string
	Type      NotificationType
	CreatedAt time.Time
	Duration  time.Duration
	closeBtn  widget.Clickable
}

// NotificationManager 通知管理器
type NotificationManager struct {
	mu            sync.RWMutex
	notifications []*Notification
	animProgress  map[string]float32
}

// NewNotificationManager 创建通知管理器
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		notifications: make([]*Notification, 0),
		animProgress:  make(map[string]float32),
	}
}

// AddNotification 添加通知
func (nm *NotificationManager) AddNotification(message string, notifType NotificationType) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	notif := &Notification{
		ID:        fmt.Sprintf("notif_%d", time.Now().UnixNano()),
		Message:   message,
		Type:      notifType,
		CreatedAt: time.Now(),
		Duration:  4 * time.Second,
	}

	nm.notifications = append(nm.notifications, notif)
	nm.animProgress[notif.ID] = 0

	// 限制最大通知数量
	if len(nm.notifications) > 5 {
		nm.notifications = nm.notifications[len(nm.notifications)-5:]
	}
}

// AddSuccess 添加成功通知
func (nm *NotificationManager) AddSuccess(message string) {
	nm.AddNotification(message, NotificationSuccess)
}

// AddError 添加错误通知
func (nm *NotificationManager) AddError(message string) {
	nm.AddNotification(message, NotificationError)
}

// AddWarning 添加警告通知
func (nm *NotificationManager) AddWarning(message string) {
	nm.AddNotification(message, NotificationWarning)
}

// AddInfo 添加信息通知
func (nm *NotificationManager) AddInfo(message string) {
	nm.AddNotification(message, NotificationInfo)
}

// RemoveNotification 移除通知
func (nm *NotificationManager) RemoveNotification(id string) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	for i, notif := range nm.notifications {
		if notif.ID == id {
			nm.notifications = append(nm.notifications[:i], nm.notifications[i+1:]...)
			delete(nm.animProgress, id)
			break
		}
	}
}

// GetActiveNotifications 获取活动通知
func (nm *NotificationManager) GetActiveNotifications() []*Notification {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	result := make([]*Notification, len(nm.notifications))
	copy(result, nm.notifications)
	return result
}

// Update 更新通知状态（检查过期）
func (nm *NotificationManager) Update(gtx layout.Context) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	now := time.Now()
	activeNotifications := make([]*Notification, 0)

	for _, notif := range nm.notifications {
		age := now.Sub(notif.CreatedAt)

		// 检查关闭按钮点击
		if notif.closeBtn.Clicked(gtx) {
			delete(nm.animProgress, notif.ID)
			continue
		}

		// 检查是否过期
		if age < notif.Duration {
			activeNotifications = append(activeNotifications, notif)

			// 更新动画进度
			fadeInDuration := 300 * time.Millisecond
			fadeOutDuration := 300 * time.Millisecond
			showDuration := notif.Duration - fadeInDuration - fadeOutDuration

			if age < fadeInDuration {
				// 淡入阶段
				nm.animProgress[notif.ID] = float32(age) / float32(fadeInDuration)
			} else if age < fadeInDuration+showDuration {
				// 显示阶段
				nm.animProgress[notif.ID] = 1.0
			} else {
				// 淡出阶段
				fadeOutProgress := float32(age-fadeInDuration-showDuration) / float32(fadeOutDuration)
				nm.animProgress[notif.ID] = 1.0 - fadeOutProgress
			}
		} else {
			delete(nm.animProgress, notif.ID)
		}
	}

	nm.notifications = activeNotifications
}

// Layout 渲染通知列表
func (nm *NotificationManager) Layout(gtx layout.Context) layout.Dimensions {
	nm.Update(gtx)

	notifications := nm.GetActiveNotifications()
	if len(notifications) == 0 {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	// 在屏幕右上角显示通知
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 右上角偏移
			return layout.Inset{
				Top:   unit.Dp(16),
				Right: unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 垂直排列通知
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.End,
				}.Layout(gtx, nm.buildNotificationList(gtx, notifications)...)
			})
		}),
	)
}

// buildNotificationList 构建通知列表
func (nm *NotificationManager) buildNotificationList(gtx layout.Context, notifications []*Notification) []layout.FlexChild {
	children := make([]layout.FlexChild, 0, len(notifications))

	for _, notif := range notifications {
		notif := notif // 捕获循环变量
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return nm.layoutNotification(gtx, notif)
		}))
	}

	return children
}

// layoutNotification 渲染单个通知
func (nm *NotificationManager) layoutNotification(gtx layout.Context, notif *Notification) layout.Dimensions {
	// 获取动画进度
	nm.mu.RLock()
	progress := nm.animProgress[notif.ID]
	nm.mu.RUnlock()

	// 应用透明度动画
	opacity := progress
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 1 {
		opacity = 1
	}

	// 应用位移动画（从右侧滑入）
	offsetX := (1 - progress) * 100

	// 通知样式
	bgColor := nm.getNotificationColor(notif.Type)
	icon := nm.getNotificationIcon(notif.Type)

	// 通知尺寸
	notifWidth := unit.Dp(320)
	notifHeight := unit.Dp(64)

	// 使用宏操作应用变换
	macro := op.Record(gtx.Ops)

	dims := layout.Stack{}.Layout(gtx,
		// 背景
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := image.Point{
				X: gtx.Dp(notifWidth),
				Y: gtx.Dp(notifHeight),
			}
			return drawNotificationBackground(gtx, bgColor, size)
		}),
		// 内容
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:  unit.Dp(16),
				Right: unit.Dp(12),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// 图标
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return nm.drawIcon(gtx, icon, notif.Type)
						})
					}),
					// 消息文本
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, 14, notif.Message)
						lbl.Color = toNRGBA(theme.DefaultTheme.Text)
						lbl.MaxLines = 2
						return lbl.Layout(gtx)
					}),
					// 关闭按钮
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return nm.layoutCloseButton(gtx, &notif.closeBtn)
						})
					}),
				)
			})
		}),
	)

	call := macro.Stop()

	// 应用变换
	op.Offset(image.Pt(int(offsetX), 0)).Add(gtx.Ops)

	// 应用透明度
	if opacity < 1 {
		paint.ColorOp{Color: color.NRGBA{A: uint8(255 * opacity)}}.Add(gtx.Ops)
	}

	call.Add(gtx.Ops)

	return dims
}

// getNotificationColor 获取通知类型对应的颜色
func (nm *NotificationManager) getNotificationColor(notifType NotificationType) color.RGBA {
	switch notifType {
	case NotificationSuccess:
		return color.RGBA{R: 40, G: 100, B: 60, A: 255}
	case NotificationError:
		return color.RGBA{R: 120, G: 40, B: 40, A: 255}
	case NotificationWarning:
		return color.RGBA{R: 140, G: 110, B: 40, A: 255}
	case NotificationInfo:
		return color.RGBA{R: 40, G: 80, B: 120, A: 255}
	default:
		return theme.DefaultTheme.Surface
	}
}

// getNotificationIcon 获取通知图标（使用Unicode字符作为简单图标）
func (nm *NotificationManager) getNotificationIcon(notifType NotificationType) string {
	switch notifType {
	case NotificationSuccess:
		return "✓"
	case NotificationError:
		return "✕"
	case NotificationWarning:
		return "!"
	case NotificationInfo:
		return "i"
	default:
		return "•"
	}
}

// drawIcon 绘制通知图标
func (nm *NotificationManager) drawIcon(gtx layout.Context, icon string, notifType NotificationType) layout.Dimensions {
	iconColor := nm.getNotificationIconColor(notifType)
	size := gtx.Dp(unit.Dp(24))

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sizePoint := image.Point{X: size, Y: size}
			return drawCircle(gtx, iconColor, sizePoint)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, 14, icon)
			lbl.Color = toNRGBA(theme.DefaultTheme.Text)
			return lbl.Layout(gtx)
		}),
	)
}

// getNotificationIconColor 获取通知图标颜色
func (nm *NotificationManager) getNotificationIconColor(notifType NotificationType) color.RGBA {
	switch notifType {
	case NotificationSuccess:
		return theme.DefaultTheme.Success
	case NotificationError:
		return theme.DefaultTheme.Error
	case NotificationWarning:
		return theme.DefaultTheme.Warning
	case NotificationInfo:
		return theme.DefaultTheme.Primary
	default:
		return theme.DefaultTheme.TextSecondary
	}
}

// layoutCloseButton 渲染关闭按钮
func (nm *NotificationManager) layoutCloseButton(gtx layout.Context, clickable *widget.Clickable) layout.Dimensions {
	return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := gtx.Dp(unit.Dp(20))
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				sizePoint := image.Point{X: size, Y: size}
				return drawCircle(gtx, theme.DefaultTheme.Border, sizePoint)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, 12, "×")
				lbl.Color = toNRGBA(theme.DefaultTheme.TextSecondary)
				return lbl.Layout(gtx)
			}),
		)
	})
}

// drawNotificationBackground 绘制通知背景
func drawNotificationBackground(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	defer clip.RRect{
		Rect: image.Rectangle{Max: size},
		SE:   8,
		SW:   8,
		NE:   8,
		NW:   8,
	}.Push(gtx.Ops).Pop()

	// 绘制背景
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// 绘制左边框（强调色）
	borderWidth := 4
	borderRect := image.Rectangle{
		Max: image.Point{X: borderWidth, Y: size.Y},
	}
	defer clip.Rect(borderRect).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(getNotificationBorderColor(c))}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}

// getNotificationBorderColor 获取通知边框颜色
func getNotificationBorderColor(bg color.RGBA) color.RGBA {
	// 根据背景色计算边框颜色（稍微亮一点）
	return color.RGBA{
		R: min(255, bg.R+40),
		G: min(255, bg.G+40),
		B: min(255, bg.B+40),
		A: 255,
	}
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}
