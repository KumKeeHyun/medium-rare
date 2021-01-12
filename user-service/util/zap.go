package util

import (
	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"go.uber.org/zap"
)

func BuildZapLogger() (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(config.App.ZapConfig.Level)); err != nil {
		return nil, err
	}

	zapCfg := zap.Config{
		OutputPaths:       config.App.ZapConfig.OutputPaths,
		DisableCaller:     !config.App.ZapConfig.EableCaller,
		DisableStacktrace: !config.App.ZapConfig.EableCaller,
		Level:             level,
		Encoding:          config.App.ZapConfig.Encoding,
		// EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}
	return zapCfg.Build()
}
