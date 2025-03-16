package utils

import (
	"SeaChat/pkg/entity"
	"image/color"

	"github.com/mojocn/base64Captcha"
)

func NewCaptcha(store base64Captcha.Store) *base64Captcha.Captcha {
	mathDriver := base64Captcha.NewDriverMath(40, 160, 5, base64Captcha.OptionShowSineLine, &color.RGBA{
			R: 254,
			G: 254,
			B: 254,
			A: 254,
		}, base64Captcha.DefaultEmbeddedFonts, []string{"wqy-microhei.ttc"})
	return base64Captcha.NewCaptcha(mathDriver, store)
}

func GenerateCaptcha(store base64Captcha.Store)(*entity.CaptchaData,error){
	if id,base64,_,err := NewCaptcha(store).Generate(); err != nil {
		return nil,err
	} else {
		return &entity.CaptchaData{
			ID: id,
			Base64: base64,
		},nil
	}
}

func VerifyCaptcha(store base64Captcha.Store, id string, value string) bool {
	return NewCaptcha(store).Verify(id, value, true)
}
