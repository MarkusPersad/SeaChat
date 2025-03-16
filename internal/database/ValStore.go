package database

import (
	"SeaChat/pkg/constants"
	"context"

	"github.com/rs/zerolog/log"
)

type ValStore interface {
	Set(id string, value string) error
	Get(id string, clear bool) string
	Verify(id, answer string, clear bool) bool
}

var ctx = context.Background()

func (s *service) Set(id string, value string) error {
	key := constants.CAPTCHA_ID + id
	return s.valkeyclient.Do(ctx, s.valkeyclient.B().Setex().Key(key).Seconds(constants.CAPTCHA_TIMEOUT).Value(value).Build()).Error()
}
func (s *service) Get(id string, clear bool) string {
	key := constants.CAPTCHA_ID + id
	result := s.valkeyclient.Do(ctx, s.valkeyclient.B().Get().Key(key).Build())
	if result.Error() != nil {
		return ""
	}
	if clear {
		e := s.valkeyclient.Do(ctx, s.valkeyclient.B().Del().Key(key).Build()).Error()
		if e != nil {
			log.Logger.Error().Err(e).Msg("valkey del error")
			return ""
		}
	}
	val, e := result.ToString()
	if e != nil {
		log.Logger.Error().Err(e).Msg("valkey get error")
		return ""
	}
	return val
}

func (s *service) Verify(id, answer string, clear bool) bool {
	val := s.Get(id, clear)
	return val == answer
}
