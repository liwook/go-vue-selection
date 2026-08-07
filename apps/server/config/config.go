package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name      string         `mapstructure:"name"`
	Mode      string         `mapstructure:"mode"`
	Port      int            `mapstructure:"port"`
	Log       LogConfig      `mapstructure:"log"`
	Postgres  PostgresConfig `mapstructure:"postgres"`
	SnowFlake SnowFlake      `mapstructure:"snowflake"`
	Static    Static         `mapstructure:"static"`
	Pprof     PprofConfig    `mapstructure:"pprof"`
}

// PprofConfig 控制独立的 pprof 调试端口。
// 采用独立端口 + 仅绑定 localhost 的方案，常驻但零额外暴露面，
// 主业务端口（port）不会暴露 /debug/pprof。
type PprofConfig struct {
	Enabled bool `mapstructure:"enabled"` // 是否开启 pprof 调试 server，默认 false
	Port    int  `mapstructure:"port"`    // 调试端口，默认 6060；仅绑定 127.0.0.1
}

type Static struct {
	Host string `mapstructure:"host"`
	Path string `mapstructure:"path"`
}

type SnowFlake struct {
	MachineID int64 `mapstructure:"machine_id"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Output     string `mapstructure:"output"` // file | stdout | both（默认 file）
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type PostgresConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	SearchPath   string `mapstructure:"search_path"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// changeHooks 保存配置热更新时回调的钩子（如日志等级热更新）。
// 采用注册机制而非直接调用具体包，是为了避免 config 反向依赖下游包造成循环引用。
var changeHooks []func(*AppConfig)

// OnChange 注册一个配置变更钩子，配置热更新（viper 监听）触发时会传入最新的 *AppConfig。
func OnChange(h func(*AppConfig)) {
	changeHooks = append(changeHooks, h)
}

func Init(configPath string) (*AppConfig, error) {
	conf := new(AppConfig)
	//viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	//viper.SetConfigType("yaml") // 指定配置文件类型（专用于从远程获取配置信息时指定配置文件类型的）
	viper.SetConfigFile(configPath)
	viper.AddConfigPath(".")                  // 指定查找文件的路径（配合相对路径使用）
	viper.AutomaticEnv()                      // 指定支持从环境变量读取配置
	viper.SetEnvPrefix("VUE_ADMIN")           // 指定环境变量 KEY 的前缀
	replacer := strings.NewReplacer(".", "_") // 替换规则，将 . 替换为 _
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err: %v\n", err)
		return nil, err
	}
	// 把读取到的配置信息，反序列化到 conf 变量中
	if err := viper.Unmarshal(conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v \n", err)
		return nil, err
	}

	// 必填校验：配置缺失应启动即报错（Fail Fast），而非带零值运行。
	if err := validate(conf); err != nil {
		fmt.Printf("config validate failed, err:%v\n", err)
		return nil, err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		slog.Info("配置文件修改了...", "file", in.Name)
		if err := viper.Unmarshal(conf); err != nil {
			slog.Error("viper.Unmarshal failed", "error", err)
			return
		}
		// 通知所有已注册的配置变更钩子（如日志等级热更新）
		for _, h := range changeHooks {
			h(conf)
		}
	})

	return conf, nil
}

// validate 校验关键必填配置，缺失即返回错误以触发 Fail Fast。
func validate(c *AppConfig) error {
	if c.Name == "" {
		return fmt.Errorf("config: app.name is required")
	}
	if c.Log.Filename == "" {
		return fmt.Errorf("config: log.filename is required")
	}
	if c.Log.Level == "" {
		return fmt.Errorf("config: log.level is required")
	}
	if c.Postgres.Host == "" {
		return fmt.Errorf("config: postgres.host is required")
	}
	if c.Postgres.User == "" {
		return fmt.Errorf("config: postgres.user is required")
	}
	if c.Postgres.DbName == "" {
		return fmt.Errorf("config: postgres.dbname is required")
	}
	return nil
}
