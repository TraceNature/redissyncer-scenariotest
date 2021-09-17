package global

import (
	"testcase/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

var (
	RSPViper  *viper.Viper
	RSPLog    *zap.Logger
	RSPConfig config.Server
	once      sync.Once
)

