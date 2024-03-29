package main

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var lvl zapcore.Level
	if err := lvl.Set(os.Getenv("LEVEL")); err != nil {
		log.Printf("cannot parse log level %s: %s", os.Getenv("LEVEL"), err)

		lvl = zapcore.WarnLevel
	}

	logFile, _ := os.Create("log.txt")

	encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	defaultCore := zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(os.Stderr)), lvl)
	fileCore := zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl)
	cores := []zapcore.Core{
		defaultCore,
		fileCore,
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	{
		logger.Debug("debug massage", zap.Time("now", time.Now()))
	}
	{
		logger.Warn("warn massage", zap.Int("data", rand.Int()), zap.Int("key", rand.Int()))
	}
	{
		logger.Info("info massage", zap.String("id", "9902777462"))
	}
	{
		logger.Error("error massage", zap.Error(errors.New("failed to load")))
	}
}
