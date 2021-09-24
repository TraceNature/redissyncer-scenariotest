package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
	"testcase/config"
)

var (
	RSPViper  *viper.Viper
	RSPLog    *zap.Logger
	RSPConfig config.Server
	once      sync.Once
)

