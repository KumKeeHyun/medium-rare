package config

import (
	"os"
	"strings"
)

type AppConfig struct {
	Address       string        `yaml:"address"`
	JWTSecret     string        `yaml:"jwtSecret"`
	ArticleConfig ServiceConfig `yaml:"articleConfig"`
	MysqlConfig   DataConfig    `yaml:"mysqlConfig"`
	ZapConfig     LogConfig     `yaml:"zapConfig"`
	KafkaConfig   KafkaConfig   `yaml:"kafkaConfig"`
}

type ServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	URL  string `yaml:"url"`
}

type DataConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type LogConfig struct {
	OutputPaths []string `yaml:"outputPaths"`
	Level       string   `yaml:"level"`
	Encoding    string   `yaml:"encoding"`
	EableCaller bool     `yaml:"enableCaller"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

var App = &AppConfig{}

func getEnvStr(env string, def string) (envValue string) {
	envValue = os.Getenv(env)
	if envValue == "" {
		envValue = def
	}
	return
}

func getEnvStrAry(env string, def []string) (envValue []string) {
	tempEnv := os.Getenv(env)
	if tempEnv == "" {
		envValue = def
		return
	}
	envValue = strings.Split(tempEnv, ",")
	return
}

func init() {
	App.Address = getEnvStr("APP_ADDR", "0.0.0.0:8084")
	App.JWTSecret = getEnvStr("APP_JWTSECRET", "kkh")

	App.ArticleConfig.Host = getEnvStr("APP_ARTICLE_HOST", "127.0.0.1")
	App.ArticleConfig.Port = getEnvStr("APP_ARTICLE_PORT", "8082")
	App.ArticleConfig.URL = getEnvStr("APP_ARTICLE_URL", "/v1/articles/list")

	App.MysqlConfig.Driver = getEnvStr("APP_DB_DRIVER", "mysql")
	App.MysqlConfig.Host = getEnvStr("APP_DB_HOST", "127.0.0.1")
	App.MysqlConfig.Port = getEnvStr("APP_DB_PORT", "3306")
	App.MysqlConfig.User = getEnvStr("APP_DB_USER", "root")
	App.MysqlConfig.Password = getEnvStr("APP_DB_PASS", "balns")

	App.ZapConfig.OutputPaths = getEnvStrAry("APP_LOG_OUTPUTS", []string{"stdout"})
	App.ZapConfig.Level = getEnvStr("APP_LOG_LEVEL", "debug")
	App.ZapConfig.Encoding = getEnvStr("APP_LOG_ENCODING", "json")
	App.ZapConfig.EableCaller = false

	App.KafkaConfig.Brokers = getEnvStrAry("APP_KAFKA_BROKERS", []string{"127.0.0.1:9092"})
}
