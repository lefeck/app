package global

import (
	"app/config"
	ut "github.com/go-playground/universal-translator"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var (
	Config = &config.Config{}
	//DB     = &gorm.DB{}
	Trans ut.Translator
	DB    *gorm.DB
	Log   *zap.Logger
)
