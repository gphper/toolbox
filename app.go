package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"

	"github.com/kbinani/screenshot"
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

func (a *App) ScreenShot() string {

	// 截取屏幕
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		log.Fatal(err)
	}

	// 将图像转换为字节数据
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Fatal(err)
	}

	// 将字节数据转换为 base64 字符串
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}
