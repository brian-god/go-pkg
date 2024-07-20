package configs

var (
	Mode EnvMode // 开发环境
)

// EnvMode 开发环境
type EnvMode string

const (
	Development EnvMode = "dev" // 开发
	Production  EnvMode = "pro" // 生产
	Prerelease  EnvMode = "pre" // 预发布
)

type Bootstrap struct {
	Server *Server `mapstructure:"server"`
	Log    *Log    `mapstructure:"log"`
	JWT    *JWT    `mapstructure:"jwt"`
}
type Server struct {
	Port       int    `mapstructure:"port"`
	RateQPS    int    `mapstructure:"rate_qps"`
	TracerPort int    `mapstructure:"tracer_port"`
	Name       string `mapstructure:"name"`
}

type Log struct {
	OutPath    string `mapstructure:"out_path"`
	FilePrefix string `mapstructure:"file_prefix"`
	Level      int64  `mapstructure:"max_size"`
	MaxSize    int64  `mapstructure:"max_size"`
	MaxBackups int64  `mapstructure:"max_backups"`
	MaxAge     int64  `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type JWT struct {
	Issuer     string `mapstructure:"issuer"`
	SigningKey string `mapstructure:"signing_key"`
}
