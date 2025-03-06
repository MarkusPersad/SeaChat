package utils

import (
	"SeaChat/pkg/structses"
	"image/color"

	"github.com/mojocn/base64Captcha"
	"github.com/rs/zerolog/log"
)

type Captcha struct {
	*base64Captcha.Captcha
}

func NewCaptcha(store base64Captcha.Store) *Captcha {
	mathDriver := base64Captcha.NewDriverMath(40, 160, 5, base64Captcha.OptionShowSineLine, &color.RGBA{
		R: 254,
		G: 254,
		B: 254,
		A: 254,
	}, base64Captcha.DefaultEmbeddedFonts, []string{"wqy-microhei.ttc"})
	return &Captcha{
		Captcha: base64Captcha.NewCaptcha(mathDriver, store),
	}
}

func(c *Captcha) GenerateBase64()(*structses.CaptchaDateBase64,error){
	if id, b64s, _, err := c.Generate(); err != nil {
		log.Logger.Error().Err(err).Msgf("Generate Captcha Error: %v", err)
		return nil,err
	} else {
		return &structses.CaptchaDateBase64{
			Id: id,
			B64s: b64s,
		},nil
	}
}
