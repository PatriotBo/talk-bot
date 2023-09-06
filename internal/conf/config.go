package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"

	"talk_bot/internal/service/openai"
)

const configPath = "./internal/conf/" // 配置文件目录

var globalConfig *Config

// DBConfig db connection config
type DBConfig struct {
	Host        string `yaml:"host"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Port        int    `yaml:"port"`
	DbName      string `yaml:"dbName"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	Timeout     int    `yaml:"timeout"` // 超时时间 单位：秒
}

// Config 配置信息
type Config struct {
	DB     DBConfig      `yaml:"db"`
	OpenAI openai.Config `yaml:"openAI"`
}

func init() {
	filename := filepath.Join(configPath, "config.yaml")
	by, err := os.ReadFile(filename)
	if err != nil {
		p, _ := os.Getwd()
		fmt.Printf("pwd:%s \n", p)
		panic(fmt.Errorf("read config err:%v", err))
	}

	c := new(Config)
	if err = yaml.Unmarshal(by, c); err != nil {
		panic(fmt.Errorf("unmarshal config err:%v", err))
	}
	globalConfig = c
}

// GetDSN 获取 db dsn
func GetDSN() string {
	cfg := globalConfig.DB
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ds&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.Timeout)
}

// GetConfig ...
func GetConfig() *Config {
	return globalConfig
}
