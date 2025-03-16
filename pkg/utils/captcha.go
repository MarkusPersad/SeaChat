package utils

import (
	"SeaChat/pkg/entity"
	"image/color"

	"github.com/mojocn/base64Captcha"
)

type SeaCaptcha struct {
	*base64Captcha.Captcha
}

func New(store base64Captcha.Store) *SeaCaptcha {
	mathDriver := base64Captcha.NewDriverMath(40, 160, 5, base64Captcha.OptionShowSineLine, &color.RGBA{
			R: 254,
			G: 254,
			B: 254,
			A: 254,
		}, base64Captcha.DefaultEmbeddedFonts, []string{"wqy-microhei.ttc"})
	return &SeaCaptcha{
		Captcha: base64Captcha.NewCaptcha(mathDriver,store),
	}
}

func (s *SeaCaptcha) Generate()(*entity.CaptchaData,error){
	if id,base64,_,err := s.Captcha.Generate(); err != nil {
		return nil,err
	} else {
		return &entity.CaptchaData{
			ID: id,
			Base64: base64,
		},nil
	}
}
