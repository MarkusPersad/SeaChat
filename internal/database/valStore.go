package database

import (
	"SeaChat/pkg/constants"
	"context"

	"github.com/rs/zerolog/log"
)

type ValStore interface {
	// 实现 Base64CaptCha 的 Store 接口
		Set(id string, value string) error

		Get(id string, clear bool) string

		Verify(id, answer string, clear bool) bool
}

func(s *service) Set(id string,value string) error {
	key := constants.CAPTCHA + id
	return s.valkey.Do(context.Background(),s.valkey.B().Setex().Key(key).Seconds(constants.CAPTCHA_TIMEOUT).Value(value).Build()).Error()
}

func(s *service)Get(id string,clear bool) string {
	key := constants.CAPTCHA + id
	result := s.valkey.Do(context.Background(),s.valkey.B().Get().Key(key).Build());
	if result.Error() != nil {
		return ""
	}
	if clear {
		if err := s.valkey.Do(context.Background(),s.valkey.B().Del().Key(key).Build()).Error(); err != nil {
			log.Logger.Error().Err(err).Msgf("failed to delete captcha key %s", key)
		}
	}
	if val,err := result.ToString();err != nil {
		return ""
	} else {
		return val
	}
}

func (s *service) Verify(id, answer string, clear bool) bool {
	val := s.Get(id, clear)
	return val == answer
}
