package pages

import (
	"cultivation-client/internal/gui/components"
	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/network"
	"cultivation-client/internal/store"
	"gioui.org/layout"
	"gioui.org/unit"
)

type RegisterPage struct {
	username   *components.InputField
	password   *components.InputField
	confirmPw  *components.InputField
	submitBtn  *components.Button
	loginBtn   *components.Button
	errorMsg   string
	loading    bool
	onRegister func()
	onLogin    func()
}

func NewRegisterPage() *RegisterPage {
	return &RegisterPage{
		username:  components.NewInputField("用户名"),
		password:  components.NewInputField("密码"),
		confirmPw: components.NewInputField("确认密码"),
		submitBtn: components.NewButton("注册"),
		loginBtn:  components.NewButton("已有账号？登录"),
		errorMsg:  "",
		loading:   false,
	}
}

func (p *RegisterPage) SetOnRegister(fn func()) {
	p.onRegister = fn
}

func (p *RegisterPage) SetOnLogin(fn func()) {
	p.onLogin = fn
}

func (p *RegisterPage) Layout(gtx layout.Context) layout.Dimensions {
	// 处理按钮点击
	if p.submitBtn.Clicked(gtx) {
		go p.handleRegister()
	}
	if p.loginBtn.Clicked(gtx) {
		if p.onLogin != nil {
			p.onLogin()
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel("注册账号")
			title.Color = theme.DefaultTheme.Primary
			title.Size = 28
			return layout.Inset{
				Top:    unit.Dp(60),
				Left:   unit.Dp(80),
				Bottom: unit.Dp(32),
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
				Top:  unit.Dp(12),
				Left: unit.Dp(80),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.password.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:  unit.Dp(12),
				Left: unit.Dp(80),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.confirmPw.Layout(gtx)
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
				Top:  unit.Dp(20),
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
				return p.loginBtn.Layout(gtx)
			})
		}),
	)
}

func (p *RegisterPage) handleRegister() {
	p.loading = true
	p.errorMsg = ""

	if p.password.Text() != p.confirmPw.Text() {
		p.errorMsg = "两次密码不一致"
		p.loading = false
		return
	}

	_, err := network.GetAPIClient().Register(
		p.username.Text(),
		p.password.Text(),
	)

	if err != nil {
		p.errorMsg = "注册失败: " + err.Error()
		p.loading = false
		return
	}

	p.loading = false

	// 连接 WebSocket
	if err := network.GetWebSocketClient().Connect(); err != nil {
		// 忽略 WebSocket 错误
	}

	// 触发注册成功回调
	if p.onRegister != nil {
		p.onRegister()
	}
}

func (p *RegisterPage) Clear() {
	p.username.SetText("")
	p.password.SetText("")
	p.confirmPw.SetText("")
	p.errorMsg = ""
}

func (p *RegisterPage) IsLoggedIn() bool {
	return store.GetAuthStore().IsLoggedIn()
}
