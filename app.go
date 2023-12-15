package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// 截屏
func (a *App) ScreenShot() string {

	// 截取屏幕
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	// 将图像转换为字节数据
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	// 将字节数据转换为 base64 字符串
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

// 生产验证码
func (a *App) Totp(data string) string {

	// 获取当前时间
	now := time.Now()
	// 生成 TOTP 实例 otpauth://totp/GitHub:gphper?secret=
	key, err := otp.NewKeyFromURL(data)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	// 生成验证码
	totpCode, _ := totp.GenerateCode(key.Secret(), now)

	return totpCode
}

// 数据存储
func (a *App) Storage(data string) int {
	path := "toolbox.data"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
		}
	}(file)

	n, err := file.Write([]byte(data))
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	return n
}

// 数据获取
func (a *App) Get() string {

	path := "toolbox.data"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
		}
	}(file)

	var allData = make([]byte, 1024)
	n, err := file.Read(allData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	return string(allData[:n])
}
