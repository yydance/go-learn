package setting

import (
	"os"

	"github.com/caarlos0/env"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	viperlib "github.com/spf13/viper"
)

type Config struct {
	Server Server `yaml:"server"`
	Log    Log    `yaml:"log"`
	Mysql  Mysql  `yaml:"mysql"`
}
type RunConfig struct {
	Env    string `env:"ENV"`
	StdOut bool   `env:"STDOUT"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Log struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	LogFile    string `yaml:"logFile"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

type Mysql struct {
	Master MysqlInfo `yaml:"master"`
	Slave  MysqlInfo `yaml:"slave"`
}

type MysqlInfo struct {
	Host     string `yaml:"host" env:"MYSQL_HOST"`
	Port     string `yaml:"port" env:"MYSQL_PORT"`
	Username string `yaml:"username" env:"MYSQL_USERNAME"`
	Password string `yaml:"password" env:"MYSQL_PASSWORD,unset"`
	Database string `yaml:"database" env:"MYSQL_DATABASE"`
}

var Viper = viperlib.New()

func InitConfig() (config *Config) {
	loadEnv()
	setConfigFile()
	var c Config
	config = c.getConf()
	return
}

func (c *Config) getConf() *Config {

	err := Viper.Unmarshal(&c)
	if err != nil {
		Logger.Panic(err)
	}

	return c
}

func setConfigFile() {
	path, err := os.Getwd()
	if err != nil {
		Logger.Panic(err)
	}
	loadFlag()
	env := Viper.GetString("env")
	Logger.Infof("当前环境：%s", env)
	Logger.Infof("当前配置文件目录: %s", path)
	Viper.AddConfigPath(path + "/configs")
	Viper.SetConfigName("config." + env)
	Viper.SetConfigType("yaml")

	if err := Viper.ReadInConfig(); err != nil {
		Logger.Panic(err)
	}
	Viper.OnConfigChange(func(in fsnotify.Event) {
		if err := Viper.ReadInConfig(); err != nil {
			Logger.Panicf("read config error: %s", err.Error())
		}
		Logger.Info("config has been changed")
	})
	Viper.WatchConfig()

}

func loadFlag() {
	var runCfg RunConfig
	pflag.StringVar(&runCfg.Env, "env", "13", "运行环境")
	pflag.BoolVar(&runCfg.StdOut, "stdout", true, "日志是否标准输出")

	var cfg Config
	pflag.StringVar(&cfg.Server.Host, "host", "127.0.0.1", "服务监听地址")
	pflag.StringVar(&cfg.Server.Port, "port", "7090", "服务监听端口")
	pflag.StringVar(&cfg.Log.LogFile, "log-file", "/var/log/demo.log", "日志文件")
	pflag.StringVar(&cfg.Log.Level, "log-level", "debug", "日志级别")
	pflag.StringVar(&cfg.Log.Format, "log-format", "text", "日志格式")

	Viper.BindPFlags(pflag.CommandLine)
	//pflag.Parse()
}

func loadEnv() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		Logger.Errorf("%+v\n", err)
	}
	var runCfg RunConfig
	if err := env.Parse(&runCfg); err != nil {
		Logger.Errorf("%+v\n", err)
	}
	//Viper.SetEnvPrefix("env")
	Viper.AutomaticEnv()
}
