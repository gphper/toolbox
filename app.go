package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"image/png"
	"io"
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

// 截屏
func (a *App) ScreenShot() (string, error) {

	// 截取屏幕
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "", err
	}

	// 将图像转换为字节数据
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "", err
	}

	// 将字节数据转换为 base64 字符串
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// 生产验证码
func (a *App) Totp(data string) (string, error) {

	// 获取当前时间
	now := time.Now()
	// 生成 TOTP 实例 otpauth://totp/GitHub:gphper?secret=
	key, err := otp.NewKeyFromURL(data)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "", err
	}
	// 生成验证码
	totpCode, err := totp.GenerateCode(key.Secret(), now)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "", err
	}

	return totpCode, nil
}

// 数据存储
func (a *App) Storage(data string) (int, error) {
	path := "toolbox.data"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return 0, err
	}
	defer file.Close()

	n, err := file.Write([]byte(data))
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return 0, err
	}

	return n, nil
}

// 数据获取
func (a *App) Get() (string, error) {

	path := "toolbox.data"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "", err
	}
	defer file.Close()

	var allData = make([]byte, 1024)
	n, err := file.Read(allData)

	if err != nil && !errors.Is(err, io.EOF) {
		runtime.LogError(a.ctx, err.Error())
		return "", err
	}

	return string(allData[:n]), nil
}
