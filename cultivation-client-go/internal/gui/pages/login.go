package pages

import (
	"cultivation-client/internal/gui/components"
	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/network"
	"cultivation-client/internal/store"
	"gioui.org/layout"
	"gioui.org/unit"
)

type LoginPage struct {
	username    *components.InputField
	password    *components.InputField
	submitBtn   *components.Button
	registerBtn *components.Button
	errorMsg    string
	loading     bool
	onLogin     func()
	onRegister  func()
}

func NewLoginPage() *LoginPage {
	return &LoginPage{
		username:    components.NewInputField("用户名"),
		password:    components.NewInputField("密码"),
		submitBtn:   components.NewButton("登录"),
		registerBtn: components.NewButton("注册账号"),
		errorMsg:    "",
		loading:     false,
	}
}

func (p *LoginPage) SetOnLogin(fn func()) {
	p.onLogin = fn
}

func (p *LoginPage) SetOnRegister(fn func()) {
	p.onRegister = fn
}

func (p *LoginPage) Layout(gtx layout.Context) layout.Dimensions {
	// 处理按钮点击
	if p.submitBtn.Clicked(gtx) {
		go p.handleLogin()
	}
	if p.registerBtn.Clicked(gtx) {
		if p.onRegister != nil {
			p.onRegister()
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel("修仙世界")
			title.Color = theme.DefaultTheme.Primary
			title.Size = 32
			return layout.Inset{
				Top:    unit.Dp(80),
				Left:   unit.Dp(80),
				Bottom: unit.Dp(40),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return title.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left: unit.Dp(80),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.username.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:  unit.Dp(16),
				Left: unit.Dp(80),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.password.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.errorMsg != "" {
				errLabel := components.NewLabel(p.errorMsg)
				errLabel.Color = theme.DefaultTheme.Error
				return layout.Inset{
					Top:  unit.Dp(8),
					Left: unit.Dp(80),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return errLabel.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:  unit.Dp(24),
				Left: unit.Dp(80),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.submitBtn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:  unit.Dp(16),
				Left: unit.Dp(80),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.registerBtn.Layout(gtx)
			})
		}),
	)
}

func (p *LoginPage) handleLogin() {
	p.loading = true
	p.errorMsg = ""

	_, err := network.GetAPIClient().Login(
		p.username.Text(),
		p.password.Text(),
	)

	if err != nil {
		p.errorMsg = "登录失败: " + err.Error()
		p.loading = false
		return
	}

	p.loading = false

	// 注册 WebSocket 消息处理器并连接
	if err := network.GetWebSocketClient().Connect(); err != nil {
		// 忽略 WebSocket 错误，继续进入主页面
	}

	// 触发登录成功回调
	if p.onLogin != nil {
		p.onLogin()
	}
}

func (p *LoginPage) Clear() {
	p.username.SetText("")
	p.password.SetText("")
	p.errorMsg = ""
}

func (p *LoginPage) IsLoggedIn() bool {
	return store.GetAuthStore().IsLoggedIn()
}
