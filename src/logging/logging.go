package logging

import (
	"community-sig/config"

	"github.com/pkg/errors"
)

var (
	Logger BaseLogger
)

func InitZap(conf *config.LogConf) error {
	if conf != nil {
		Logger = NewZapLogger(conf)
	} else {
		Logger = &ConsoleLogger{level: InfoLevel}
	}

	if Logger == nil {
		return errors.New("InitZap error")
	}

	return nil
}
