package generatedata

import "testcase/global"

type TargetType int

const (
	TargettypeSingle  TargetType = 0
	TargettypeCluster TargetType = 1
)

func (tt TargetType) String() string {
	switch tt {
	case TargettypeSingle:
		return "single"
	case TargettypeCluster:
		return "cluster"
	default:
		return ""
	}
}

type BigKey struct {
	KeySuffixLen uint `mapstructure:"keysuffixlen" json:"keysuffixlen" yaml:"keysuffixlen"`
	Length uint `mapstructure:"length" json:"length" yaml:"length"`
	ValueSize uint `mapstructure:"valuesize" json:"valuesize" yaml:"valuesize"`
	Expire int64 `mapstructure:"expire" json:"expire" yaml:"expire"`
	Duration uint `mapstructure:"duaration" json:"duaration" yaml:"duaration"`
	DataGenInterval uint `mapstructure:"datageninterval" json:"datageninterval" yaml:"datageninterval"`
}

type RandKey struct {
	KeySuffixLen uint `mapstructure:"keysuffixlen" json:"keysuffixlen" yaml:"keysuffixlen"`
	ValueSize uint `mapstructure:"valuesize" json:"valuesize" yaml:"valuesize"`
	Expire int64 `mapstructure:"expire" json:"expire" yaml:"expire"`
	Duration uint `mapstructure:"duaration" json:"duaration" yaml:"duaration"`
	DataGenInterval uint `mapstructure:"datageninterval" json:"datageninterval" yaml:"datageninterval"`
	Threads int `mapstructure:"threads" json:"threads" yaml:"threads"`
}

type GenData struct {
	TargetType TargetType `mapstructure:"type" json:"type" yaml:"type"`
	Addr       []string     `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password   string     `mapstructure:"password" json:"password" yaml:"password"`
	DB int `mapstructure:"db" json:"db" yaml:"db"`
	BigKey	BigKey `mapstructure:"bigkey" json:"bigkey" yaml:"bigkey"`
	RandKey RandKey `mapstructure:"randkey" json:"randkey" yaml:"randkey"`

}

func(gd *GenData) Exec(){
	global.RSPLog.Sugar().Info("GenData execute")
}