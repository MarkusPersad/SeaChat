package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cilium/lumberjack/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	maxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
	maxBackups, _ := strconv.Atoi(os.Getenv("LOG_MAX_BACKUPS"))
	maxAge, _ := strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
	isCompressed, _ := strconv.ParseBool(os.Getenv("LOG_COMPRESSED"))
	lv := os.Getenv("LOG_LEVEL")
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("./logs/log-%s.log", currentDate)
	hook := lumberjack.Logger{
		Filename:   logFileName,
		MaxAge:     maxAge,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		Compress:   isCompressed,
	}
	switch strings.ToLower(lv) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}
	zerolog.TimeFieldFormat = os.Getenv("LOG_TIME_FORMAT")
	log.Logger = log.Logger.With().Caller().Stack().Logger()
	if strings.ToLower(lv) == "debug" {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: os.Getenv("LOG_TIME_FORMAT"),
		}
		multi := zerolog.MultiLevelWriter(&hook, consoleWriter)
		log.Logger = log.Output(multi)
	} else {
		log.Logger = log.Output(&hook)
	}
}
