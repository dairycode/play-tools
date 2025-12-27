package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

var AppConfig *Config

// Init 从 YAML 文件加载配置
func Init() (err error) {
	// 读取配置文件
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return
	}

	// 解析 YAML
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return
	}

	// 校验并填充默认值
	setDefaultValues(config)

	AppConfig = config
	log.Println("Configuration loaded from config.yaml")
	return
}

// setDefaultValues 校验配置并设置默认值
func setDefaultValues(config *Config) {
	// 数据库配置默认值
	if config.Database.User == "" {
		config.Database.User = "root"
		config.Database.Password = "123456"
		log.Println("Database.User not set, using default: root")
	}
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
		log.Println("Database.Host not set, using default: localhost")
	}
	if config.Database.Port == "" {
		config.Database.Port = "3306"
		log.Println("Database.Port not set, using default: 3306")
	}
	if config.Database.Name == "" {
		config.Database.Name = "play_tools"
		log.Println("Database.Name not set, using default: play_tools")
	}

	// 服务器配置默认值
	if config.Server.Port == "" {
		config.Server.Port = "8080"
		log.Println("Server.Port not set, using default: 8080")
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
		log.Println("Server.Host not set, using default: 0.0.0.0")
	}
}
