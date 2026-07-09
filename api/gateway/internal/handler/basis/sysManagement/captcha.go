package sysManagement

import (
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var captchaFontFace font.Face

func init() {
	fnt, err := opentype.Parse(goregular.TTF)
	if err == nil {
		face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
			Size:    28,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err == nil {
			captchaFontFace = face
		}
	}
}

type CaptchaHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCaptchaHandler(svcCtx *svc.ServiceContext) *CaptchaHandler {
	return &CaptchaHandler{
		svcCtx: svcCtx,
	}
}

func (h *CaptchaHandler) GenerateCaptcha(w http.ResponseWriter, r *http.Request) {
	keyLen := h.svcCtx.Config.Captcha.KeyLong
	if keyLen <= 0 {
		keyLen = 6
	}
	imgWidth := h.svcCtx.Config.Captcha.ImgWidth
	if imgWidth <= 0 {
		imgWidth = 240
	}
	imgHeight := h.svcCtx.Config.Captcha.ImgHeight
	if imgHeight <= 0 {
		imgHeight = 80
	}

	code := generateRandomCode(keyLen)
	id := generateCaptchaID()

	GetCaptchaStore().Set(id, code, 5*time.Minute)

	img := drawCaptchaImage(code, imgWidth, imgHeight)

	var buf strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	png.Encode(encoder, img)
	encoder.Close()

	api.OkWithDetailed(w, map[string]interface{}{
		"captcha_id":     id,
		"pic_path":       "data:image/png;base64," + buf.String(),
		"captcha_length": keyLen,
	}, "captcha")
}

func generateRandomCode(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

func generateCaptchaID() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func drawCaptchaImage(code string, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{245, 245, 245, 255}}, image.Point{}, draw.Src)

	for i := 0; i < 30; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		noiseColor := color.RGBA{
			uint8(rand.Intn(200)),
			uint8(rand.Intn(200)),
			uint8(rand.Intn(200)),
			100,
		}
		drawLine(img, x1, y1, x2, y2, noiseColor)
	}

	face := captchaFontFace
	textColor := color.RGBA{50, 50, 50, 255}

	charSpacing := width / (len(code) + 1)
	textY := height/2 + 14

	for i, ch := range code {
		offsetY := rand.Intn(10) - 5
		rotation := float64(rand.Intn(20) - 10)
		_ = rotation

		point := fixed.Point26_6{
			X: fixed.Int26_6((charSpacing*(i+1) - 5) * 64),
			Y: fixed.Int26_6((textY + offsetY) * 64),
		}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(textColor),
			Face: face,
			Dot:  point,
		}
		d.DrawString(string(ch))
	}

	for i := 0; i < 500; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		dotColor := color.RGBA{
			uint8(rand.Intn(200)),
			uint8(rand.Intn(200)),
			uint8(rand.Intn(200)),
			80,
		}
		img.Set(x, y, dotColor)
	}

	return img
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}
	if steps == 0 {
		img.Set(x1, y1, col)
		return
	}
	for i := 0; i <= steps; i++ {
		x := x1 + dx*i/steps
		y := y1 + dy*i/steps
		if x >= 0 && x < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
			img.Set(x, y, col)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
